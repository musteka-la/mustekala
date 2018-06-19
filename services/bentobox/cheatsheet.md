## Cheatsheet

### Get the last block in the canonical chain

	curl -X POST -H 'Content-Type: application/json' --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":83}' http://tigress.musteka.la:8545

### Get the block body. i.e. block header + txs + ommers

	curl -X POST -H 'Content-Type: application/json' --data '{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["0x1b4", true],"id":1}' http://tigress.musteka.la:8545

### Get the info of one specific tx (you don't really need that, see above)

	curl -H 'Content-Type: application/json' -X POST --data '{"jsonrpc":"2.0","method":"eth_getTransactionByHash","params":["0x3f9d0c7c39f9dad2a41254af6ca3e531be06db9739cbf9d9343c52e592267adc"],"id":1}' http://tigress.musteka.la:8545

### Get the transaction receipt

	curl -H 'Content-Type: application/json' -X POST --data '{"jsonrpc":"2.0","method":"eth_getTransactionReceipt","params":["0x3f9d0c7c39f9dad2a41254af6ca3e531be06db9739cbf9d9343c52e592267adc"],"id":1}' http://tigress.musteka.la:8545