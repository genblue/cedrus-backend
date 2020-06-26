pragma solidity ^0.5.12;

import "@openzeppelin/contracts-ethereum-package/contracts/token/ERC20/ERC20Detailed.sol";
import "@openzeppelin/contracts-ethereum-package/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts-ethereum-package/contracts/access/roles/WhitelistAdminRole.sol";
import "@openzeppelin/contracts-ethereum-package/contracts/math/SafeMath.sol";
import "@openzeppelin/upgrades/contracts/Initializable.sol";
import "@openzeppelin/contracts-ethereum-package/contracts/GSN/GSNRecipient.sol";

interface DaiToken {
    function transferFrom(address src, address dst, uint wad) external returns (bool);
}

contract CedrusToken is Initializable, ERC20Detailed, ERC20, GSNRecipient, WhitelistAdminRole {

    using SafeMath for uint256;
    // Store minting limits for addresses
    mapping (address => uint256) minters;

    // Dai token address
    address public daiAddress;

    function initialize(string memory name, string memory symbol, uint8 decimals) initializer public {
        ERC20Detailed.initialize(name, symbol, decimals);
        WhitelistAdminRole.initialize(msg.sender);
        GSNRecipient.initialize();
    }

    modifier isMintingValid(address minter, uint256 mintingAmount) {
        uint256 currentLimit = minters[minter];

        // Checks whether the minter is approved for their minting amount
        require(currentLimit >= mintingAmount, 'insufficient mint limit');

        // Reduce current limit of minter by the minting amount
        minters[minter] = currentLimit.sub(mintingAmount);
        _;
    }

    function acceptRelayedCall(
        address relay,
        address from,
        bytes calldata encodedFunction,
        uint256 transactionFee,
        uint256 gasPrice,
        uint256 gasLimit,
        uint256 nonce,
        bytes calldata approvalData,
        uint256 maxPossibleCharge
    ) external view returns (uint256, bytes memory) {
        return _approveRelayedCall();
    }

    function _preRelayedCall(bytes memory context) internal returns (bytes32) {
        // solhint-disable-previous-line no-empty-blocks
    }

    function _postRelayedCall(bytes memory context, bool success, uint256 actualCharge, bytes32 preRetVal) internal {
        // solhint-disable-previous-line no-empty-blocks
    }


    function getDaiAddress() public view returns (address) {
        return daiAddress;
    }

    // Allows the token owner to set dai token address
    function updateDaiAddress(address daiAddress_) public onlyWhitelistAdmin {
        daiAddress = daiAddress_;
    }

    // Allows the token owner to set minting limits
    function updateMinter(address minter, uint256 newLimit) public onlyWhitelistAdmin {
        minters[minter] = newLimit;
    }

    // Add a new admin address
    function addAdmin(address account) public onlyWhitelistAdmin {
        super.addWhitelistAdmin(account);
    }

    // Allows an approved minter to mint new Cedar Coins for an address after
    // validating their donations off-chain
    function mintCedarCoin(address claimAddress, uint256 cedarAmount)
        public 
        isMintingValid(_msgSender(), cedarAmount)
    {
        _mint(claimAddress, cedarAmount);
    }

    // Mint new Cedar Coins for an address after transferring sender's ERC20 tokens to the minter address
    function exchangeTokens(address minter, uint256 cedarAmount)
        public
        isMintingValid(minter, cedarAmount)
    {
        // Need first to be approved by the token contract
        if (!DaiToken(daiAddress).transferFrom(_msgSender(), minter, cedarAmount)) revert("Could not transfer remote tokens");

        _mint(_msgSender(), cedarAmount);
    }
}
