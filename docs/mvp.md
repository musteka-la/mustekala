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