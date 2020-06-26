const { accounts, contract } = require('@openzeppelin/test-environment');
const { BN, expectEvent, expectRevert } = require('@openzeppelin/test-helpers');
const { expect } = require('chai');

const CedrusToken = contract.fromArtifact('CedrusToken');
const DAIToken = contract.fromArtifact('CedrusToken');

describe('CedrusToken', function () {
    const [ owner, alice, bobMinter, claimer ] = accounts;

    beforeEach(async function () {
        this.cedrusContract = await CedrusToken.new({ from: owner });
        this.cedrusContract.initialize("Cedar", "CDR", 18, { from: owner });
        this.daiContract = await DAIToken.new({ from: owner });
        this.daiContract.initialize("DAI", "DAI", 18, { from: owner });
    });

    it('fail minting when minter not approved', async function () {
        // Given
        const claimAmount = 10;

        // When Then
        await expectRevert(
            this.cedrusContract.mintCedarCoin(claimer, claimAmount, { from: bobMinter }),
            'insufficient mint limit');
    });

    it('mint Cedar coin from a new minter address', async function () {
        // Given
        const claimAmount = 100;
        const expectedAmount = new BN(claimAmount);
        await this.cedrusContract.updateMinter(bobMinter, 100000, { from: owner });

        // When
        const receiptMint = await this.cedrusContract.mintCedarCoin(claimer, claimAmount, { from: bobMinter });
        expectEvent(receiptMint, 'Transfer', { from: "0x0000000000000000000000000000000000000000", to: claimer, value: expectedAmount });

        // Then
        expect((await this.cedrusContract.balanceOf(claimer)).toString()).to.equal('100');
    });

    it('fail minting Cedar coin when a minter has a mint limit', async function () {
        // Given
        const claimAmount = 100;
        const mintLimit = 10;
        await this.cedrusContract.updateMinter(bobMinter, mintLimit, { from: owner });

        // When Then
        await expectRevert(
            this.cedrusContract.mintCedarCoin(claimer, claimAmount, { from: bobMinter }),
            'insufficient mint limit');
    });

    it('fail purchasing Cedar when no DAI tokens on the DAI contract', async function () {
        // Given
        const claimAmount = 100;
        const mintLimit = 1000;
        await this.cedrusContract.updateMinter(bobMinter, mintLimit, { from: owner });
        await this.cedrusContract.updateDaiAddress(this.daiContract.address, { from: owner });

        // When Then
        await expectRevert(
            this.cedrusContract.exchangeTokens(bobMinter, claimAmount, { from: alice }),
            'ERC20: transfer amount exceeds balance');
    });

    it('fail purchasing Cedar when cedar token address not approved on DAI contract', async function () {
        // Given
        const claimAmount = 100;
        const mintLimit = 1000;
        await this.cedrusContract.updateMinter(bobMinter, mintLimit, { from: owner });
        await this.cedrusContract.updateDaiAddress(this.daiContract.address, { from: owner });
        // We use the Cedar contract as a ERC20 contract
        await this.daiContract.updateMinter(owner, mintLimit, { from: owner });
        await this.daiContract.mintCedarCoin(alice, claimAmount, { from: owner });

        // When Then
        await expectRevert(
            this.cedrusContract.exchangeTokens(bobMinter, claimAmount, { from: alice }),
            'ERC20: transfer amount exceeds allowance');
    });

    it('purchase Cedar with DAI tokens', async function () {
        // Given
        const claimAmount = 100;
        const mintLimit = 1000;
        const expectedAmount = new BN(claimAmount);
        await this.cedrusContract.updateMinter(bobMinter, mintLimit, { from: owner });
        await this.cedrusContract.updateDaiAddress(this.daiContract.address, { from: owner });
        // We use the Cedar contract as a ERC20 contract
        await this.daiContract.updateMinter(owner, mintLimit, { from: owner });
        await this.daiContract.mintCedarCoin(alice, claimAmount, { from: owner });
        await this.daiContract.approve(this.cedrusContract.address, claimAmount, { from: alice });

        // When
        const receiptMint = await this.cedrusContract.exchangeTokens(bobMinter, claimAmount, { from: alice });

        // Then
        // 1. We expect a DAI Transfer event
        expect(receiptMint.logs[0].event.toString()).to.equal('Transfer');
        expect(receiptMint.logs[0].args.from.toString()).to.equal(alice);
        expect(receiptMint.logs[0].args.to.toString()).to.equal(bobMinter);
        expect(receiptMint.logs[0].args.value.toString()).to.equal("100");

        // 2. We expect a DAI Approval event (with 0 approval remaining)
        expect(receiptMint.logs[1].event.toString()).to.equal('Approval');
        expect(receiptMint.logs[1].args.owner.toString()).to.equal(alice);
        expect(receiptMint.logs[1].args.spender.toString()).to.equal(this.cedrusContract.address);
        expect(receiptMint.logs[1].args.value.toString()).to.equal("0");

        // 3. We expect a CEDAR Transfer event
        expect(receiptMint.logs[2].event.toString()).to.equal('Transfer');
        expect(receiptMint.logs[2].args.from.toString()).to.equal('0x0000000000000000000000000000000000000000');
        expect(receiptMint.logs[2].args.to.toString()).to.equal(alice);
        expect(receiptMint.logs[2].args.value.toString()).to.equal("100");

        expect((await this.cedrusContract.balanceOf(alice)).toString()).to.equal('100');
        expect((await this.daiContract.balanceOf(alice)).toString()).to.equal('0');
        expect((await this.daiContract.balanceOf(bobMinter)).toString()).to.equal('100');
    });
});
