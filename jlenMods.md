## JL's storage branch
JL created a `storage_rent` branch in May 2018. [changes here](https://github.com/optimalbrew/rskj/compare/master...optimalbrew:storage_rent#) 

This may be prior to the unitrie, i.e. separate trie for account and storaeg (as in eth).

### Block executor
rskj-core/src/main/java/co/rsk/core/bc/BlockExecutor.java

- imported address class, but did not use it. 

### Account State
rskj-core/src/main/java/org/ethereum/core/AccountState.java

The current version of AccountState (post unitrie) is different. Changes noted here just to understand logic at that point.
- items encoded in RLP: `nonce, balance, stateroot, codehash, stateFlags, lastModificationTime `. But in current version, its just `nonce, balance,` and `stateFlags` (no storage root or codehash).

- line 76: introduced `lastModificationTime` for an account (cache last mod time)
- line 125-135: add new method `set` and `get` methods for `LastModificationDate`
- changes to node state marked via `setDirty(true)` so these can be saved at end of block execution


### Repository
 rskj-core/src/main/java/org/ethereum/core/Repository.java 

**aside:** repository is an *interface* class to keep track of the state in memory. methods to create accounts, retrieve account info, dump state of current repository, store all state changes to actual DB (`commit`), or `rollback` undo changes since some snapshot, etc.  

- modified transfer() so TX that has a payment of 0 value `return`s immediately. No need to change state (balances).
- no changes to storage. 
**aside** but the sender will still be charged gas for making the call, yeah?


 rskj-core/src/main/java/org/ethereum/db/RepositoryTrack.java

- L241 method to check is an account is cached and a mod to skip making redundant state changes (e.g. balance delta is 0) 


### Transaction structure and summary
rskj-core/src/main/java/org/ethereum/core/Transaction.java

TX structure: how transactions should be created, encoded, decoded. 

- new var `rentGasLimit` with `get` and `set` methods, and also `RLP.encode/decode` for new TX structure to account for rent.
- add new possibilities for TX calling structure and TX `create` method (new method signature using rent).  

rskj-core/src/main/java/org/ethereum/core/TransactionExecutionSummary.java
 - routine mods, mostly following the logic for regular gas


### Transaction execution
rskj-core/src/main/java/org/ethereum/core/TransactionExecutionSummary.java 

- new vars `storageprice`, but also `codePrice` and `accontStatePrice` (separate pricing)
- a bunch of changes cloasely tracking the logic for regular gas
- L661: `rentGasCalculator(addr)`  very close to rskip113 logic
- L649: `setCorrectTimes(addr)` timestamp for state mod, with +6 months for new nodes
- L432: `payRent()` uses above (timestamp and rent calculator), generates OOG exception for rent.
- and many other changes, along somewhat predictable lines.


### Contracts
rskj-core/src/main/java/org/ethereum/db/ContractDetailsCacheImpl.java

- L199: changes to two methods `getStorageKeys()` and `getStorageSize()`. These should be different now cos unitrie.


### VM
rskj-core/src/main/java/org/ethereum/vm/VM.java
- L648 looks odd (repeated)
- L781: mods to `EXTCODECOPY` and `EXTCODESIZE`

rskj-core/src/main/java/org/ethereum/vm/program/Program.java
- large file 2K .. mods seem simple, but needs study for proper context. 

rskj-core/src/main/java/org/ethereum/vm/program/ProgramResult.java
- mostly mods to keep track of accounts that have been modified or newly created.

### Misc

rskj-core/src/main/java/org/ethereum/util/TimeUtils.java
- create a constant for 6 months in seconds (advanced rent payment).

### Tests
A bunch of tests to go with all the mods.
