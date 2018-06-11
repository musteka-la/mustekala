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

**Kitsunet** (The fox network) is the mesh of MetaMask browser peers.

This layer is concerned with discovering and connecting peers (via *libp2p*), as
well as consuming, storing and sending *Ethereum-IPLD* data among these peers.

The optimization on the transport of information (as well as its storage), taking
advantage of the merkle trie data structure is domain of the **Layer 4** of this
architecture.

#### What do we know so far?

* Discovering: By using the js-ipfs client, we can leverage its DHT to discover
  **libp2p** peers. This approach is insuficient, as we want, `a)` Find
  *kitsunet* nodes and `b)` Have a granular approach to discovery, i.e.
  privilege the finding of bridges and hubs, as well as other *kitsunet* nodes
  running elements our client must be interested (See **Layer 4**, subsets).

* Connecting: As our *kitsunet* clients will be running in browsers, it is
  rather difficult to assign an arbitrary transport port to our service (as,
  for example, the port `8545`). The compromise is found on using proxy servers
  which our clients will connect to. An scaling protocol to facilitate the
  connections among peers must be found, as well as an *easy to package*
  solution to offer the users, so they can run their own **Hubs**. The goal:
  Enable the *Kitsunet* clients to find their closest *Hub* to be part of the
  mesh.

* Consuming: While IPLD is not packaged in IPFS clients *out-of-the-box*, that
  is, if you want to parse *ethereum-ipld* nodes, you need to have the plugin
  installed in your client; The IPFS clients are able to deal with these
  elements as *DAG Nodes* (*directed acyclic graph*). In other words, we can
  ask any IPFS client for an ethereum ipld node, by using its hash. We will
  realize soon, though, that this approach is elemental and not efficient, as
  in order to get, for example, the balance of an ethereum account, starting
  from a certain block header, you need to perform a non-trivial number of
  traversal queries (6 to 8 as of 2018.06.11). It is imperative to work on
  an overlay able to boost the exchange protocol of *DAG nodes*, taking
  advantage of the blockchain merkle tree data structure. (See **Layer 4**).

* Storing: This element is fundamental for the **Layer 4** of the architecture,
  as each **kitsunet** client will store a number of subsets of the ethereum
  state, as well as co-selector indexes.

* Sending: This is the ability to deliver requested data to other peers, such
  as, the latest block header (*pubsub*), ethereum state (*layer 4*) and
  co-selector indexes. The ability to perform this feature represents of the
  departure of light clients as mere consumers of information from hubs.

#### MVP Features

* Simple *js-ipfs* client connecting via *libp2p* with a bridge and fetching
  account and storage data. (aka *Naive mode*).
* Rudimentary subscription to new block header.
* Updating the **balance of an account**, on each new block header.
* Updating the information of a particular smart contract by pulling its
  storage trie from a **Bridge**. (No traversing, but at once).
* **EVM Call** to extract information from this smart contract.
* Partial adoption into MetaMask clients
  * Telemetry monitoring.
  * Adopt a "I have the subset of state XYZ" flag.
  * Discovery of peers with similar "subset of state XYZ".
  * Send transactions to the **Bridges**
* Simulation of the network

#### Beyond MVP

* Scale **kitsunet** to 1MM+ peers
  * Signalling servers
  * Rendevous servers
  * **Emphasis** on containerized docker solutions to boost infrastructure.
    and foment decentralization.
* *PubSub* system for New Block Headers.
* Transaction broadcasting and relaying throughout **kitsunet**.
* Make modules available (*js*, *go*, *python*) to enable any program to
  leverage kitsunet.
* **GOLD**: Miners connecting to **kitsunet** listening for transactions.

--------------------------------------------------------------------------------

### Layer 4: Content Routing System

#### What is this?

The promise of **IPLD** (which stands for *InterPlanetary Linked Data*), is not
only make clients be able to lookup and retrieve nodes, and gain knowledge on
other nodes linking to them. Also, it is implied that new functionalities will
be developed, such as **transformations**.

Now, without aditional features, simple operations, like for instance, the
traversal for an account take (as of 2018.06.11) 6 to 8 travels from the client
to a source of data (*Bridge*, *Hub* or another client). Not to mention that
it may be necessary for certain smart contracts, to obtain a good part of its
storage trie.

On the other hand, to pursue the ideal of decentralization, and to
achieve the dream state of avoiding peers to hold huge ammount of bytes so they
can work properly (as of 2018.06.11 the ethereum state is in the
10-20GB mark), *kitsunet* needs to work as a **distributed storage of the
ethereum data**.

We want to accomplish this by developing a custom **Content Routing System**
(**CRS**) able to function as an overlay to the *DAG Node* exchange protocol,
and optimized to the use case of the Ethereum State.

#### What do we know so far?

* [IPLD Use Case for Ethereum Light Client](https://github.com/ipld/ipld/issues/29)

* Features of the **Content Routing System** (**CRS**)
  * Divide the state into **subsets** (avoiding the word *shard*),
  which will be small, redundant, well spaced and useful to each peer.

  * Each peer of **kitsunet**, which, by the way can be not only a browser peer,
  but a **Hub**, will maintain a number of these **subsets**, consisting on an
  organized number of ethereum state and storage trie nodes, and will update
  their elements as the *Block Header* of the *Canonical Chain* goes changing.

  * A peer will maintain, ideally, the *subsets* containing relevant data to its
  operation, plus a couple of discrete *subsets* to ensure redundancy and
  availability of the whole system.

  * Maintainers of *subsets* will require at each block header update an
  **index** of a certain *subset* (ex: Index for subset `0x1a56`), to known
  peers maintaining such *subset*. This query is ultimately located to a
  **Bridge** peer in charge of maintaining this data by synchronizing it from
  *devp2p*.

  * An **index** consists on a list containing the first two bytes of each hash
  belonging the *subset* of the state. Thus, for example, the index of the
  subset `0x1a56` can be the list `[0x4505, 0xa5ac, 0x34ab...]`.

  * The peer computes the received **index** against the stored **index** for
  that subset, noting the differences (deltas) between them. These **deltas**
  enable the peer to prepare the list of **needed nodes**.

  * As the peer, by this process, can't know the **needed nodes** hashes at
  full length; It needs to request them by their *relative reference*.
  This is a reference including the *id* of the **subset** and their position.
  For example, it will require "For the subset `0x1a56`, I need the nodes `2a`
  and `35`, `36`", meaning that it needs the *ath child of the 2nd child* of
  the head of this **subset** as well as the
  `5th and 6th children of the 3rd child`.

  * The **kitsunet CRS** will locate the peers providing these nodes and send
  them to the requester.

  * It won't be discarded the capacity of sending full **subsets** on demand,
  to facilitate fast synchronizations.

#### MVP and Beyond Features

* Rudimentary Proof of Concept of description above.
* Research
  * "*Hot Caching*" of **Bridges** and, **Hubs**
    * This is, a *hub* receives queries of a certain **subset**. If resources
    are available, this *hub* will start maintaining this **subset** to
    increase local redundancy and availability.
  * *PubSub* features must be builtin into the **kitsunet CRS**
  * Optimizations on *Co-selector indexes* to be outside **Bridges** and living
  in **Hubs** and browser peers.
* Relevant metrics
  * Redundancy
  * Availability
  * Well Spacing
  * Clustering (i.e. "How far of my neighbourhood I have to go to get data?")
  * Data Exchange Rate

#### The future!

* COMPLETE Substitution of Ethereum JSON/RPC use.
  * Updated Data Retrieval
  * EVM Calls
  * Transactions Broadcasting and Relaying.

* Nice `localhost` WebApps to complement the **kitsunet** mesh
  * State of the network
  * Block Explorer
  * Log Explorer

--------------------------------------------------------------------------------
