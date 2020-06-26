#!/usr/bin/env sh

set -euxo pipefail

# Available Accounts
# ==================
# (0) 0x9F27999FC0230cD177a8E28F90D85adBF90589C9 (Owner & Upgrader)
# (1) 0x0569790EE0343DB5ce92e9D4E566544a3e752448 (Minter)
# (2) 0xEe7Ca29E972D3738A4cDAfF537e838522CE4BAc0 (Alice)

owner="0x9F27999FC0230cD177a8E28F90D85adBF90589C9"
minter="0x0569790EE0343DB5ce92e9D4E566544a3e752448"
aliceUser="0xEe7Ca29E972D3738A4cDAfF537e838522CE4BAc0"

# 1. Deploy of CedrusToken contract
echo "1. Deploy of CedrusToken contract\n"
echo "Please copy/pasted the deployed CedrusToken as you will be asked to enter it.\n"

./node_modules/.bin/oz deploy \
 CedrusToken \
 --kind upgradeable \
 --from ${owner} \
 --network rinkeby

# 2. Add new minter
echo "2. Add new minter\n"
./node_modules/.bin/oz send-tx \
 --network rinkeby \
 --from ${owner} \
 --method updateMinter \
 --args "${minter}, 12000000000000000000"

# 3. Mint some coins to Alice
echo "3. Mint some coins to Alice\n"
./node_modules/.bin/oz send-tx \
 --network rinkeby \
 --from ${minter} \
 --method mintCedarCoin \
 --args "${aliceUser}, 2000000000000000000"

read -p "Enter the Cedrus deployed address : " cedrusAddress

# 4. Send the gasless transaction using the cedar-lib
echo "4. Send the gasless transaction using the cedar-lib\n"
node -r esm ./test/scripts/spend.js "${cedrusAddress}" "${minter}" "${aliceUser}"
