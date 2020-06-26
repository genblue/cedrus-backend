const CedrusLib = require('../../../cedar-lib');
const { projectId } = require('../../secrets.json');

console.log("Cedrus contract address", process.argv[2])
console.log("Cedrus minter address", process.argv[3])
console.log("Alice address", process.argv[4])

const contractAddress = process.argv[2];
const minterAddress = process.argv[3];
const aliceAddress = process.argv[4];

let cedrusLib = new CedrusLib.default(
    `https://rinkeby.infura.io/v3/${projectId}`,
    contractAddress,
    minterAddress,
    { privateKey: "0x742314a00b599f23a71642cc0a95a72497a75de6c3e731146330af0f2c6a7233" });

let response = cedrusLib.spend(aliceAddress, 1)
    .then(response => {
        console.log(response)
    })
    .catch(error => console.error(error))

console.log(response)
console.log("Successfully sent a gasless transaction !")

