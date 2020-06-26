const { projectId, mnemonic } = require('./secrets.json');
const { mnemonicProd } = require('./.mnemonic.json');
const HDWalletProvider = require('@truffle/hdwallet-provider');

module.exports = {
  networks: {
    development: {
      protocol: 'http',
      host: 'localhost',
      port: 8545,
      gas: 90000000,
      gasPrice: 5e9,
      networkId: '*',
    },
    rinkeby: {
      provider: () => new HDWalletProvider(
          mnemonic, `https://rinkeby.infura.io/v3/${projectId}`
      ),
      networkId: 4,
      gasPrice: 11e9
    },
    mainnet: {
      provider: () => new HDWalletProvider(
          mnemonicProd, `https://mainnet.infura.io/v3/${projectId}`
      ),
      networkId: 1,
      gasPrice: 32e9
    }
  },
};
