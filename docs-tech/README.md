# Services

```shell
.
├── Makefile
├── README.md
├── api-docs            # Swagger documentation
│   ├── cedrusservice
│   └── sethservice
├── assets
│   ├── email           # Email template
│   └── static          # Frontend
├── bin                 # Binaries
├── build
│   └── package         # Docker materials
├── cmd                 # Entrypoints
├── docker-compose.yaml
├── docs
├── internal            # Internal packages of services
├── package.json
├── pkg                 # Shared packages
├── test                # Test assets
└── vendor              # Third-party assets
```

## `cedrusservice`

The service host the frontend API endpoints, it is also a proxy to the `sethservice` `/transfer` endpoint.  
API Swagger can be found at the path: `http://localhost:8000/documentation/swagger/index.html#/`

Environment variables (with examples):
    
    API_PORT: 8000
    DB_URI: mongodb://mongo:27017
    SETH_SERVICE_URL: http://sethservice:8002/api/v1/transfer
    DB_NAME: cedrus

## `emailservice`

This service execute as a cronjob, fetch unsent emails in the database and send it.

Environment variables (with examples):

    DB_URI: mongodb://mongo:27017
    DB_NAME: cedrus

## `sethservice`

This service sends transactions to the blockchain, and is not meant to be exposed directly to internet.  
API Swagger can be found at the path: `http://localhost:8002/documentation/swagger/index.html#/`

Environment variables (with examples):

    API_PORT: 8002
    DB_URI: mongodb://mongo:27017
    DB_NAME: cedrus
    TOKEN_CONTRACT_ADDRESS: "0xde0B295669a9FD93d5F28D9Ec85E40f4cb697BAe"
    KEYSTORE_WALLET_PASSWORD: ""
    SETH_CHAIN: "ethlive"
    ETH_KEYSTORE: "/home/sethservice/secrets/ethereum"
    ETH_PASSWORD: "/home/sethservice/secrets/ethereum/password"
    ETH_FROM: "0xaD16D6D10e6Acf06C6A17Bd85cFb9f1D5467C644"
    ADMINISTRATOR_EMAIL: dummyaddres@gmail.com
    ADMINISTRATOR_NAME: Rick The Admin
