# source notes

## Trie
Trie [package](https://github.com/rsksmart/rskj/tree/49ab8dee512fb18c745e7044d51e0d57ee32da18/rskj-core/src/main/java/co/rsk/trie)

In the implementations, several references to `trie` should be read as `trieNode`, and not the full tree data structure.

* Each node has an optional associated value (a byte array)
* A node is referenced via a *key* (a byte array). See class `NodeReference`
* A node can be serialized to/from a message (a byte array)
* A node has a hash (keccak256 of its serialization)
* A node is **immutable**: to add/change a **value** or **key**, a **new node is created** (hence the need for *mutabletrie*)
* An empty node has no subnodes and a null value

The `trie` (node) constructor has these fields for each node (many of which are `final`)
- `hash`, always 32 bytes. Used as the key to retrieve and store nodes in store.
- `store` the kv data store. 
- (an optional) `value, valueLength` (useful for `EXTCODESIZE`), and `valuehash`
- `left` and `right` children
- `childrenSize` (size of entire subtree including the node itself)
- `sharedPath`  immutable slice of a node's key (of class `TrieKeySlice`)

Many other attributes and some methods: `arity=2`, `encoded` (removed before saving), `fromMessage` for serialization from store/  A `trieSize` method to count number of nodes in subtree. This is in constract to `NodeReference.nodeSize`   

Since a trie (node) is not mutable, `put` operations lead to creation of new nodes (trie node objects).


## NodeReference

Constructor for the `Nodereference` class needs
- a (trie) `store`
- optionally a node and a hash (initialized to `null` if no node is passed) 

Methods:
- `getNode()`: if the ndoe is present (cached) return it, if node is null but we have the hash, then retrive it from `store`. If node and has are both null, then create an `empty` node ie `Nodereference(null, null, null)`. 
- similarly a `gethash()`
- **note**: `nodeSize`returns the sum of3 `trie` methods
    - `getValueLength` + `getChildrenSize` + `getMessageLength`

**concern:** `trie.getchldrenSize()` calls `NodeReference.referenceSize()`. Both have warnings and should not be called from outside. Reference size itseif makes a call to getchildrenSize (via `nodeSize`).. so this has a recursive structure.

For rent computations, it is better to disregard information on childrensize..  and use the sum of `valuelength` and `messagelength` only. Since rent is node based and not account based, that's how it should be anyway.
- for terminal nodes (`trie.isTerminal()`), this is no issue at all 
- however if rent is to be computed for *intermediate* nodes, then not including childrensize avoids accounting errors (such as double counting). 
- RSKIP113 says there is no rent for intermediate nodes. SL communicated that transactions should be charged for intermediate nodes that are loaded during regular node lookups. This is to serve as a disincentive for attackes via lookups for non-existnent accounts. However, there may be alternative ways to penalize such wild-goose chases (without causing must distress for honest mistakes or bugs in routine TXs.    



## TrieStoreImpl
Think node store, rather than trie store. 

This implementation stores and retrieve nodes (byte array serialization) by a node's `hash`.

Has a `save()` method to recursively save all unsaved nodes.


## MutableTrie
An interface. 

It is implemented in 
`MutableTrieImpl` in [co.rsk.db](https://github.com/rsksmart/rskj/blob/master/rskj-core/src/main/java/co/rsk/db/MutableTrieImpl.java)
- has `getStorageKeys(addr)`and `StorageKeysIterator`

and `MutableTrieCache` in [co.rsk.db](https://github.com/rsksmart/rskj/blob/master/rskj-core/src/main/java/co/rsk/db/MutableTrieCache.java)
- a single cache to mark both *changed* and *deleted* elements.
- 


# Repository

https://github.com/rsksmart/rskj/blob/master/rskj-core/src/main/java/org/ethereum/core/Repository.java#L39

`AccountState createAccount(RskAddress addr);` // only for accounts, not for storage

* This method creates an account, but is DOES NOT create a contract.
* To create a contract, internally the account node is extended with a root node for storage. 
* To avoid creating the root node for storage each time a storage cell is added, we pre-create the storage node when we know the account will become a contract. This is done in `setupContract()`
* Note that we can't use the length or existence of the code node for this, because a contract's code can be empty!

`void setupContract(RskAddress addr);` //this is used for contracts



