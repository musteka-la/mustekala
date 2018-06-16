## Trie Generator

Creates an Ethereum secure state trie for a given number of accounts using
Redis DB as backend.

### Quick start

From the root directory of this repository

	make trie-generator && ./build/bin/trie-generator --accounts 100

Will create a secure state trie with 100 accounts, returning its root.

Make sure to store that root somewhere, if you want to operate on the trie afterwards.

### Database structure

There is no database structure, keys are stored as binaries in redis.