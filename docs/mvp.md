## MVP features

The following are the planned features for an MVP.

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
    and run their own, with minimal maintainance

