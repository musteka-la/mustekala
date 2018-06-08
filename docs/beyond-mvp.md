## Beyond MVP

This is the wishlist, the vision for the project

* Reduce the distance from the **devp2p** network to the bucket
  * Optimize go-ethereum to insert the data directly into the **layer 2**
  * Optimize parity to insert the data directly into the **layer 2**
  * Merge both layers 1 and 2. Have an ethereum client connect to **libp2p**
    as well, making all its stored keys available at once
* Pubsub system for new block headers
* Storage of co-selector indexes
  * With a location address system
* The network protocol is not specified in the yellow paper
  * Meaning that the "winning" network is the one where the minting block
    consensus happens

* Implement ipfs-cluster as **Bridge**
* Scalable, clusterizable database (SQL DBMS) as IPFS datastore
* Efficient *pub/sub* system to notify on new blocks
  * And indexing of subsets of the state (see **layer 4**)
* Ethereum transactions should be facilitated by the **Hubs** and **Bridges**
* Pruning of data service (highly configurable)
* Maintainance of co-señector indexes, such as
  * getLogs support (ugh)
  * CHT (block number -> block hash)
  * block tracker / block head syncer
  * tx block references (tx -> block)
  * block expansion (block -> txs)
* Multi-blockchain support
  * As long as it uses content-addressed data

