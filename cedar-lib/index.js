import Web3 from 'web3' ;
import { GSNProvider} from '@openzeppelin/gsn-provider';
import CedrusContract from '../contracts/build/contracts/CedrusToken.json';
import DAIContract from '../contracts/build/contracts/IERC20.json';

export default class CedrusLib {
    constructor(web3Provider, cedrusContractAddress, daiContractAddress, minterAddress, options = null) {
        this.web3 = new Web3(new GSNProvider(web3Provider, {
            useGSN: true,
            signKey: (options && options.privateKey) ? options.privateKey : null
        }));
        this.cedrusContractAddress = cedrusContractAddress.toLowerCase();
        this.daiContractAddress = daiContractAddress.toLowerCase();
        this.minterAddress = minterAddress.toLowerCase();

        this.cedarContract = new this.web3.eth.Contract(CedrusContract.abi, this.cedrusContractAddress, {data: CedrusContract.bytecode});
        this.DAIContract = new this.web3.eth.Contract(DAIContract.abi, this.daiContractAddress);
        // Set provider without GSN wrapping
        this.DAIContract.setProvider(web3Provider)
    }

    async verifyDAIAddress() {
        // Make sure the DAI contract address is the same as the one defined in the Cedar smart-contract
        const daiAddressInCedarContract = await this.cedarContract.methods.getDaiAddress().call();
        if(daiAddressInCedarContract.toLowerCase() !== this.daiContractAddress) {
            throw new Error(`DAI contract address in library and smart-contract are different.`);
        }
    }

    async spend(fromAddress, amount) {
        // Sends the transaction via GSN
        let txReceipt;
        try {
            const amountWithDecimals = amount * 1e18;
            txReceipt = await this.cedarContract.methods
                .transfer(this.minterAddress, amountWithDecimals.toString())
                .send({ from: fromAddress });
        } catch (e) {
            throw new Error(`Insufficient funds (${e.message})`);
        }

        return txReceipt;
    }

    async watchCedar() {
        const tokenSymbol = 'CDR';
        const tokenDecimals = 18;
        const tokenImage = 'http://cedarcoin.org/images/icon.svg';

        this.web3.givenProvider.send(
            {
                method: 'wallet_watchAsset',
                params: {
                    type: 'ERC20',
                    options: {
                        address: this.cedrusContractAddress,
                        symbol: tokenSymbol,
                        decimals: tokenDecimals,
                        image: tokenImage,
                    },
                },
                id: Math.round(Math.random() * 100000),
            },
            (err, added) => {
                if (added) {
                    return added;
                } else {
                    throw new Error("Error adding token: " + err);
                }
            }
        );
    }

    async approveDAITransfer(fromAddress, amount) {
        const amountWithDecimals = amount * 1e18;
        // Sends the transaction without GSN
        return this.DAIContract.methods
            .approve(this.cedrusContractAddress, amountWithDecimals.toString())
            .send({ from: fromAddress });
    }

    async buyCedarWithDAI(fromAddress, amount) {
        const amountWithDecimals = amount * 1e18;
        // Sends the transaction via GSN
        return this.cedarContract.methods
            .exchangeTokens(this.minterAddress, amountWithDecimals.toString())
            .send({ from: fromAddress });
    }
}
