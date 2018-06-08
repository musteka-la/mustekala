## Architecture

Consists in 4 layers:

* Layer 1: **devp2p** data sources
* Layer 2: IPFS Bridges and Hubs
* Layer 3: Kitsunet Peers
* Layer 4: Content Routing System

### Layer 1: **devp2p** data sources

#### What is this?

The data of the Ethereum blockchain is copied (one or another way)
in their peers. As these peers synchronize, the state trie, accounts and smart
contracts update. Also, the clients keep copies of block headers, transactions,
transaction receipts and some useful indexes.

This layer comprises all the processes incurred to take the ethereum data, and
pass it into the following layer of IPFS bridges and hubs.

#### What do we know so far?

Obtaining this data is not a straightforward process. The main implementations,
namely, go-ethereum and parity use single-threaded databases for optimum
performance: levelDB and RocksDB respectively.

The current state of the art (as of 2018.06.08) allows to obtain the needed
data in two ways:

* By having a client discover the peers of the network andconnecting to them,
and running the **devp2p** protocols *eth/62* and *eth/63* (aka "fast") and
asking for the data needed. While this is the _de-facto_ process, it has been
observed with a very low rate of success, both finding what we can call
"good nodes" as the velocity of perform the actual synchronization.

* Is important to mention that the storage size of all the ethereum states
(i.e. an "Archive Node") is (as 2017.06.08) in the vicinity og 1TB. A single
block state is in the range of 10GB to 20GB. Which makes unsustainable for
the casual user to maintain.

* Moreover, while each main client offers a flavor of light client, there are
a number of caveats, namely, the chance of finding a synchronized run willing
to accept this light client connection.

* Also. As the synchronization process consist on processing the transactions
indicated on each block, in order to be able to update the stored stated, the
user finds herself constantly with a high usage of resources. Making this
process expensive for consumer grade devices.

* The option, so far, is to rely on third party services, offering an RPC
gateway to their synchronized nodes. Besides the situation of centralization
it is desided to avoid, we find that the offered RPC methods only restricts to
a subset of the data. i.e. You can't just ask for a complete traversal of the
state trie at the second children of the block B.

* Leaving aside for a moment the concrete issue of having the ability and
availability of certain subset of data or another from any point of the network
(the way we think a real p2p network should work), there is the issue of the
**EVM CALLS**: Generally speaking, to be able to get information from a smart
contract, its code should be executed, and its complete storage trie has to be
available. As the ethereum data is locked-in inside the clients, the only way
to fully execute this smart contract properly is by sending this call to a
synchronized peer. The latter won't scale well, should a service receive tens or
hundred of such calls.

#### MVP Features

* Pull the ethereum data from synchronized go-ethereum and parity clients
  * block headers
  * ommers
  * transaction
  * transaction receipts
  * state trie nodes
    * (in a limited way, as, "not the whole state per block")
  * storage trie nodes
    * (see above)
  * logs
* Dump the obtained information into an optimized bucket
  * See **layer 2**
* Monitoring thoughput of pulled data

#### Beyond MVP

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

--------------------------------------------------------------------------------

### Layer 2: IPFS Bridges and Hubs

#### What is this?

The obtained data must be make available to **libp2p**, we use a combination of
IPFS, with the ability to link ethereum hashes in-protocol with IPLD.

#### What do we know so far?

* Work has been done in enabling the IPFS clients to deal with Ethereum
  * https://github.com/ipld/js-ipld
  * https://github.com/ipfs/go-ipld-eth
* We need to determine how well will the protocol scale to millions, maybe
  billions of keys. It is up to the **layer 4** to determine a Content Routing
  System to deal with the complexity of the state
* In other words, we may need to build an overlay protocol, optimized to the
  ethereum data structure (namely, the patricia merkle tree)

#### MVP Features

* Run a reverse proxied IPFS client as a **Bridge**
  * This is considered at an early stage "*the node 0*"
  * However, the code of this bridge, as well as its containerization will be
    released, so anybody can set up an equal bridge
* Monitoring throughput of inserted data
* No pruning of data is contemplated
* Run a small network of **Hubs**
  * **Hubs** are IPFS nodes that pull data from other IPFS peers (as opposed to
    pulling it from **devp2p**), under a number of logics to be determined
  * **Hubs** are the backbone of **kitsunet**, everybody can get their container
    and run their own, with minimal maintainance

#### Beyond MVP

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

--------------------------------------------------------------------------------

### Layer 3: Kitsunet Peers

#### What is this?

(TODO)

#### What do we know so far?

(TODO)

#### MVP Features

(TODO)

#### Beyond MVP

(TODO)

--------------------------------------------------------------------------------

### Layer 4: Content Routing System

#### What is this?

(TODO)

#### What do we know so far?

(TODO)

#### MVP Features

(TODO)

#### Beyond MVP

(TODO)

--------------------------------------------------------------------------------
