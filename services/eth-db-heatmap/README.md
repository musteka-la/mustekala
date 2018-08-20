
## ETHEREUM DATABASE HEATMAP

### What is it

Traverses the state trie of a given block, printing information of interest:
(needs a disconnected (i.e. "cold") geth levelDB source of data)

  - Number of state trie nodes
    - types: branch, extension or leaves
    - storage used on state trie nodes

  - Number of ethereum accounts vs smart contracts

  - Smart Contracts, Number of storage trie nodes
    - types: branch, extension or leaves
    - storage used on storage trie nodes

### How it works

- By taking "slices" of the trie
  - A slice is a subtrie, starting from a head arbitrary from the root
  - For example, the slice 0x1a04, is the head found in the path 1-a-0-4
  - Of course, the root slice is the one which its head is the very root
- We cheat and after we compute each slice, we store them into the DB
  - So next time we ask for it, is a O(1)...
- The stored content is
  - Stem of the slice
  - Hashes and Values of each element
  - Which hashes are at the bottom of the slice
    - (kind of the "leaves" of this slice)
  - Sum of bytes of the values of the elements
  - Stats of the fetch
    - Retrieval from DB time
    - How many RLP branches, extension, leaves
	- While we are processing the info, the program will output these stats

Once you have the slices processed, take it to the next level!
Build a nice viz.

### Internals

#### Input params per slice

Input params are the way we identify the slices

- state-root
- head
- depth