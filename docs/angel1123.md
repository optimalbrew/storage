# AL on trie
https://github.com/rsksmart/rskj/pull/1123

AL strongly believes that the Trie and the repository should be decoupled.


Several interfaces, classes work with Tries. For example..

## Repository
Bunch of interfaces
- `Repository` interface in [ethereum.core](https://github.com/rsksmart/rskj/blob/master/rskj-core/src/main/java/org/ethereum/core/Repository.java)
- extends the `RepositorySnapshot` interface in [rsk.db](https://github.com/rsksmart/rskj/blob/master/rskj-core/src/main/java/co/rsk/db/RepositorySnapshot.java)
- which in turn, extends the `AccountInformationProvider` interface in [rsk.core.bc](https://github.com/rsksmart/rskj/blob/master/rskj-core/src/main/java/co/rsk/core/bc/AccountInformationProvider.java)


## Trie
`Trie` class implemented in [co.rsk.trie](https://github.com/rsksmart/rskj/blob/master/rskj-core/src/main/java/co/rsk/trie/Trie.java)
- interface `MutableTrie` also in [co.rsk.trie](https://github.com/rsksmart/rskj/blob/master/rskj-core/src/main/java/co/rsk/trie/MutableTrie.java) mutates the top node of a parent trie.. so this is prior to Unitrie?
- the interface is implemented in `MutableTrieImpl` in [co.rsk.db](https://github.com/rsksmart/rskj/blob/master/rskj-core/src/main/java/co/rsk/db/MutableTrieImpl.java)
- and `MutableTrieCache` also in [co.rsk.db](https://github.com/rsksmart/rskj/blob/master/rskj-core/src/main/java/co/rsk/db/MutableTrieCache.java)

Angels PR 1123:  The `save()` method in any of these can be used to save the Trie. This can be problematic. So the PR removes `getTrie()`, `getRoot()` and `save()` methods from (some of?) them  

To do this, it creates two new class implementations: 1. `TopRepository` to manage Trie changes in a transparent manner, and 2. `RepositoryTrack` handles creation of child repository.


## SL
Not opposed to improving code and performance. But worries about impact on existing code and plans for storage rent and parallel transaction execution.
- because these changes will prevent node access (through the existing classes)
- These changes may be more in line with RSKIP61 (where the lastrent time was at the account level). However, this may conflict with rskip113, where the tracking for rent is at the node level (timestamps).
- per account tracking can be problematic for TX paralellism

Per SL, an improvement to delete AccountStatePrefix check from MutableTrie and move it to the Repository 

Per AL: a trie wrapper can be added to new classes to allow per node access, so old proposals can work just fine. 



## Misc

AccountInformationProvider:
- 