# Repository and Trie
https://github.com/rsksmart/rskj/pull/1123

Several interfaces and classes work with Tries. Many of them have `save` methods, which can be problematic. AL thinks the Trie and the repository should be decoupled. SL agrees with motivation, but is concerned about side effects and consequences for storage rent and parallel execution.

## Repository
Basic interfaces
- `AccountInformationProvider` interface in [rsk.core.bc](https://github.com/rsksmart/rskj/blob/master/rskj-core/src/main/java/co/rsk/core/bc/AccountInformationProvider.java)
  - given `addr`: return `nonce, balance, iscontract, storagekeys` and `keyCount`, use `keys` to retrieve storage `values` (`dataword` or raw `bytes`)
  - **Implementation:** [rsk.core.bc.PendingState](https://github.com/rsksmart/rskj/blob/8bb0e406e4ac82bbd4a7140e76e3d2319eb913d2/rskj-core/src/main/java/co/rsk/core/bc/PendingState.java)

- `RepositorySnapshot` in [rsk.db](https://github.com/rsksmart/rskj/blob/master/rskj-core/src/main/java/co/rsk/db/RepositorySnapshot.java) extends `AccountInformationProvider`
  - the snapshot interface has a collection of *read-only* methods. `isExist(addr)`, `getroot` (storage root of data repository), `getAccountsKeys`, `getCodeHash`, `getCodeHash`, `getAccountState`. Interestingly it also has a `startTracking` method to creates a new child repository (of type Repository) to track changes.

- `Repository` interface in [ethereum.core](https://github.com/rsksmart/rskj/blob/master/rskj-core/src/main/java/org/ethereum/core/Repository.java) extends `RepositorySnapshot`
  - `getTrie()`, `createAccount(addr)` (non contract, storage node not needed), `setupContract(addr)` (create storage nodes as well), `delete`, `hibernate`, `increaseNonce`, `setNonce`, `saveCode`, `addStorageRow` (value in dataword), `addStorageBytes` (value in bytes), `addBalance`, `transfer`, `commit` store changes to repository in actual database, `rollback`, `save`, `updateAccountState`.
  - **Implementation:** [eth.vm.program.storage](https://github.com/rsksmart/rskj/blob/f442bf2e04de10cb08efdab3289b57fda689dc4d/rskj-core/src/main/java/org/ethereum/vm/program/Storage.java) and [eth.db.mutableRepository](https://github.com/rsksmart/rskj/blob/f442bf2e04de10cb08efdab3289b57fda689dc4d/rskj-core/src/main/java/org/ethereum/db/MutableRepository.java)
    

## Trie
- `Trie` class implemented in [co.rsk.trie](https://github.com/rsksmart/rskj/blob/master/rskj-core/src/main/java/co/rsk/trie/Trie.java)
  - a central data structure, lots of methods, 1500 sloc
  - The full constructor for a Trie object `value` (bytes), `left, right` child nodes, `TrieStore`, `sharedPath` (using a class `trieKeySlice` ), `valuelength`, `valuehash`, `chlidrenSize` (entire subtree, not just left and right), and a method `checkValueLength()`.
  - methods for node: `get`, `put`, `delete`, `getLeft/Right`, `isTerminal`, `getValue`, `getValueLength`, `checkValueLength`, `getChildrenSize`, methods for de/serialization `toMessage` and `fromMessage` (different ones for Orchid and Wasabi), 
  - methods for Trie: `put`, `isEmptyTrie`, tree traversal iterators (in/pre-post order) 

- interface `TrieStore`: abstract Trie methods `save, flush, retrieve, dispose`
  - [implementation](https://github.com/rsksmart/rskj/blob/master/rskj-core/src/main/java/co/rsk/trie/TrieStoreImpl.java)
- `NodeReference` class:
  - `getNode, getHash, nodesize`. This class is the only hit for **nodesize** in the entire rskj repo.
  - `nodesize` takes a trie as input and returns size of children + external value length + getMsgLength
  - From a rent perspective, need to take care of what to charge for here.. if at all (if a node has children it is an intermediate node).
  - This computation of nodesize would make sense in the ethereum world where each account's storage has its own trie (and we can use that as one basis for rent computations).
- interface `MutableTrie` also in [co.rsk.trie](https://github.com/rsksmart/rskj/blob/master/rskj-core/src/main/java/co/rsk/trie/MutableTrie.java) 
  - mutates the top node of a parent trie
  - methods `getValueLength`, used for optimizing EXTCODESIZE and `getvaluehash` for EXTCODEHASH. 
  - the interface is implemented in `MutableTrieImpl` in [co.rsk.db](https://github.com/rsksmart/rskj/blob/master/rskj-core/src/main/java/co/rsk/db/MutableTrieImpl.java)
   and `MutableTrieCache` in [co.rsk.db](https://github.com/rsksmart/rskj/blob/master/rskj-core/src/main/java/co/rsk/db/MutableTrieCache.java)


### PR 1123  
The `save()` method (in many of them) can be used to save the Trie. This can be problematic. So the PR removes `getTrie()`, `getRoot()` and `save()` methods from (some of?) them  

To do this, it creates two new class implementations: 1. `TopRepository` to manage Trie changes in a transparent manner, and 2. `RepositoryTrack` handles creation of child repository.


### Discussion
SL: Not opposed to improving code and performance. But worries about impact on existing code and plans for storage rent and parallel transaction execution.
- because these changes will prevent node access (through the existing classes)
- These changes may be more in line with RSKIP61 (where the lastrent time was at the account level). However, this may conflict with rskip113, where the tracking for rent is at the node level (timestamps).
- per account tracking can be problematic for TX paralellism

Per SL, an improvement to delete AccountStatePrefix check from MutableTrie and move it to the Repository 

Per AL: a trie wrapper can be added to new classes to allow per node access, so old proposals can work just fine. 

