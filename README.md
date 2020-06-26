# Cedrus

## Services

This project is the frontend & backend of the Cedrus Reference Design and it handles the following functionality
- Endpoint for static files used by the website
- REST API for creating new CedarCoin claims (`cedrusservice`)
- Sendgrid integration for sending email instructions on how to claim a CedarCoin (`emailservice`)
- REST API that wraps `seth` command line utility which allows CedarCoin to be claimed (`sethservice`)

For detailled explanation of the services, see [services](docs/README.md).

## Contracts

All `Ethereum` contracts are located in the `contracts/` folder.

We use [OpenZeppelin CLI](https://docs.openzeppelin.com/cli/2.7/) to manage contracts dependencies, deploy and upgrade contracts.

# Components
- Mongo DB (3.4.23)
- Sendgrid
- Seth
- Golang (1.12)
- NPM

## Install

> Prerequisite:  
> - Install `go`

```shell
$> make deps
```

## Build 

```shell
$> make build
```

## Test

```shell
$> make test
```

## Run locally

> ### Prerequisites
> 1 - Install:
> - `docker-compose`
> - `docker`  
>
> 2 - Create a `.env` file at the root location of the repository, and add the `SENDGRID_API_KEY` environment variable :
> ```shell
> SENDGRID_API_KEY=<YOUR_API_KEY>
> ```
> 
> Make sure to have outbound connection not blocked by your local network.
>

Then run on first time (might takes minutes to complete):
```shell
$> make build-docker
```
It will build `docker` containers of `mongodb`, `cedrusservice`, `emailservice` and `sethservice`, see `docker-compose.yaml`.  
Warning: the `sethservice` container build may takes minutes. 

Then, run the containers with:
```shell
$> make run
```
Note that `emailservice` will be restarted every 10 seconds to simulate a CRON job.

> Execute `docker ps` to show the running containers

To stop the containers, run:
```shell
$> make stop
```

## API Documentation

API Swagger can be found at the path: `http://localhost:8000/documentation/swagger/index.html#/`

To generate documentation static files, run
```shell
$> make docs
```

## Deploy to Pivotal Clound foundry

> Prerequisites:  
> 1. Make sure you have Cloud Foundy CLI [installed](https://docs.cloudfoundry.org/cf-cli/install-go-cli.html).  
> 2. Set the environments variables :
> - `REGISTRY`: Docker registry url (i.e `dockerhub-account/my-image`) 
> - `USERNAME`: Docker registry username
> - `PASSWORD`: Docker registry password

Run:
1. `make docker-images`
2. `make docker-tag`
3. `make docker-push`
4. `cf login` and select the right space whish is `cedar`
5. `make cf-deploy`

## Update minter address

To replace `minter` keystore file:
1. Get private key of `minter` from mnemonic: https://iancoleman.io/bip39/
2. `geth account import privatekey` without the `0x`
