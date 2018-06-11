### break down of rpc methods and node requirements (updated)

```
client
  eth_protocolVersion
  eth_syncing
  eth_coinbase
  eth_mining
  eth_hashrate
  eth_gasPrice
  eth_accounts
  eth_sign
  eth_sendTransaction
  eth_sendRawTransaction
  eth_newFilter
  eth_newBlockFilter
  eth_newPendingTransactionFilter
  eth_uninstallFilter
  eth_getFilterChanges
  eth_getFilterLogs

client broadcast
  eth_sendRawTransaction

ipfs
  eth_getBlockByHash
  eth_getUncleByBlockHashAndIndex

ipns/block syncing (head tracking)
  eth_blockNumber

ipld:selectors (selectors/transforms)
  eth_getBlockTransactionCountByHash
  eth_getUncleCountByBlockHash

index:txToBlock (coselector)
  eth_getTransactionReceipt
  eth_getTransactionByHash
  eth_getTransactionByBlockHashAndIndex

log query?? / geth bloomFilterTrie / index:logToTx
  eth_getLogs (+index:logToTx)

lazy vm (slow) / remote vm + proofs (faster)
  eth_call
  eth_estimateGas

index:CHT
  eth_call
  eth_estimateGas
  eth_getBalance
  eth_getStorageAt
  eth_getTransactionCount
  eth_getCode
  eth_getBlockByNumber
  eth_getTransactionByBlockNumberAndIndex (+index:txToBlock)
  eth_getBlockTransactionCountByNumber (+ipld:selectors)
  eth_getUncleCountByBlockNumber (+ipld:selectors)
  eth_getUncleByBlockNumberAndIndex
  eth_getLogs (+log query)
```