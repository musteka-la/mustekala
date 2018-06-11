## MVP features

This is the aggregation of the MVP points written down in architecture.md

(Features include, when possible, the github issue associated to them)

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
* Monitoring throughput of pulled data

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
    and run their own, with minimal maintainance.

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

* Rudimentary Proof of Concept of the Content Routing System
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
