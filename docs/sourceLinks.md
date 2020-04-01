# source links

things of note (but in no particular order)


https://github.com/rsksmart/rskj/blob/master/rskj-core/src/main/java/org/ethereum/core/Repository.java#L39

`AccountState createAccount(RskAddress addr);` // only for accounts, not for storage

* This method creates an account, but is DOES NOT create a contract.
* To create a contract, internally the account node is extended with a root node for storage. 
* To avoid creating the root node for storage each time a storage cell is added, we pre-create the storage node when we know the account will become a contract. This is done in `setupContract()`
* Note that we can't use the length or existence of the code node for this, because a contract's code can be empty!

`void setupContract(RskAddress addr);` //this is used for contracts



