## MUSTEKALA

MetaMask Light Client Development.

In one line: Ethereum on IPFS.

### Quick Links

* [Documentation of the project](https://github.com/MetaMask/mustekala/issues/4)
* [Architecture](https://github.com/MetaMask/mustekala/blob/master/docs/architecture.md)

### Services

#### Development Requirements

* Go 1.9.2
  * (See [This issue](https://github.com/ethereum/go-ethereum/issues/15752#issuecomment-354271572))
* [Redis](https://redis.io/)
  * Version 4.0.8
* Have Geth in your `$GOPATH`. (`go get github.com/ethereum/go-ethereum`)
  * Check that you have at least version `v.1.8.2` (revision .`b8b9f7f44`)

#### DevP2P Node Scrapper

Connects as a peer of the devp2p network, storing information of the visited peers into the DB.

[Go to the service README](services/devp2p-node-scrapper/README.md)
