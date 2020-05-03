# Notes on classes

## RepositorySnapshot (co.rsk.db)
- interface: read only methods of a repository
	- extends AccountInfoProvider
		- getStorageRoot(), getAccountKeys, getCodelength, getHash, getAccountState, 
		- create child repository for tracking changes		

## repository (org.eth.core)
- interface, extends the read only methods of RepositorySnapshot.. so write/create methods as well
	but still an interface.. no fields.. and no method implementatios.
	-getTrie()
	-createAccount(), setupContract(), delete(addr)
	-setNonce(), saveCode(), addStorageRow/Bytes(), addBalance()
	-Tracking changes: comit(), rollback(), save()

	## implementations?
	- org.eth.vm.prog.Storage (also implements program listener)
	- org.eth.db.MutableRepository 

# MutableRepository (org.ethereum.db)
- the basic repository is just an interface (a collection of some methods, with no fields).
- the implementation brings in the fields
	- mutableTrie
	- trieKeyMapper
- instantiated with mutableaTrieImpl(triestore, trie)
- and the repository methods are them implemented... 
	-storage via trie.put, getStorageKeys, storageValues, 
	- updateAccountState, 

# org.ethereum.VM
- constructor fields: 
	- vmconfig
	- precompiledcontracts
- execution fields: program, stack, opCode, OldMemSize, gasCost 
- methods:
	* math ops, bit shifts, logical ops, hash, 	
	* calcMemGas(oldMemSize, newMemSize, CopySize) //maxmemsize is 30bits =1GB, quadratic costs
	* spendgas(), and spendOpCodeGas()
	* opd.. methods for all OpCodes
	All opCode methods implemented as doOPCODE: 
		* e.g. doSSTORE L1272, doGas L1374 (remaining gas), 
		* doCreate L1409 (followed by doCreate2) 

	and follow a pattern

	```
	protected void do_opcode(){
		spendOpCodeGas() //called within a `computegas` block in some cases
		//execute the stack program
		..stack operations as needed..
		program.step() //done with step, increment the counter, for next OP
	} 

	```

- the last opCode is doSuicide.. :) .. then we have helper methods
	* executeOpCode a large switch that matches methods to specific OpCode
	* steps() the VM takes to execute program: -> study this.
	* play()


# org.eth.VMHook .. this is the interface..
	public interface VMHook {
    		void startPlay(Program program);
    		void step(Program program, OpCode opcode);
    		void stopPlay(Program program);
	}

# Program classes (within org.ethereum.vm.program)
+ invoke (interface and classes to invoke a program, a transfer, or suicide)
	- createProgramInvocation: from TX or from contract calls.. 	
	- many of these classes need to be updated for rent.
+ Listerner (events from program execution?, onstoragePut, onstorageClear)
	- these haven't been touched in 3 years.
- Memory
- Program
- Stack
- Storage
- ProgramResult

## Memory.Java
- implements programListenerAware
- memory read and write ops related to running a program
- details may not be reqd for rent implementation.

## Stack.java
- also implements programlistenerAware 
- for pop, push, swap


## Program.Java (in org.ethereum.vm.program)

Fields:
- final/constants: 
	* MAX-Memory (set at 1GB, 1<<30)
	* Transaction, ProgramInvoke, ProgramListener and Tracer
	* Stack, Memory, Storage, returnBuffer, programResult, programTrace
	* ops, byte[] of the seq of Ops 
	* precompiled, vmconfig, deletedAccountsInBlock 
- Not final:
	* pc (prog count? #ops in prog), lastOp, (bool) stopped, startAddr, jumpDestSet 
	* rskOwnerAddress
-methods: lots
	* getOp(pc), gerCurrentOP(), setLastOp(), StackPush(dataword)/pushOne/pushZero(), stackPop(), stackClear(),
	* get/setpc(), stop(), setHreturn(), verifyStackOverflow(), verifyStackSize()
	* getMemSize(), memorySave() which uses memory.write, or memory.extend(), memLoad uses memory.read, 
	* createContract L451 (detailed), create/create2() which just call the detailed version.
	* executeCode(),	
	* spendGas(), spendAllGas(), refundGas(), futureRefundGas(), refundRemainingGas(), getProgramResult(),
	* storageSave() L930, getOwnerAddress, getOriginAddr(), getBalance(), storageLoad (SSTORE)
	* getCode(), which is different from getCodeAt, getCodeLengthAt, getCodeHashAt,  
	* suicide()  

	
## Storage (org.ethereum.vm.program)
- also implements repository.. via fields
	- repository
	- addr
- initialized via
	- ProgramInvoke.getOwnerAddress() and programInvoke.getRepository()
	- then the bunch of repository methods to udpate trie, accountInfo

## ProgramResult.Java
- Fields include: 
	* gasUsed=0, and boolean revert, futureRefund (would need one for rent as well?)
- Has collections (sets and maps) that use Dataword. Apparently Java sets and Maps cannot disinguish dulticates if they are bytearrays! Hence, using dataword type rather than byte[]
	* Map<dataword, key> codechanges -> track changes in contract code
	* Set<dataword> deleteAccounts
	* Lists: internalTXs 
-Methods:
	* spendGas() // gasUsed = GasCost.add(gasUsed, gas);
	* setRevert(), isRevert()
	* refundGas() // gasUsed = GasCost.subtract(gasUsed, gas);
	* setters and getters: getGasUsed(), getCodeChanges(), getDeleteAccounts(), add/getFutureRefund() RENT!
	and many more.. such as merge(programResult), clearAllFieldsonException(), ..




# Collect thoughts on Gas Spending: where computed, where collected, reimbursed.


- TX executor
	* makes calls to spendGas() .. 2 calls using Program.spendGas() and 2 directly using programResult.spendGas()
	* does much more though.. it deduts gasLliit initially.. and keeps track of gas to be refunded..
	* 

- Program.Java
	* L590: getProgramResult()// where storage costs are deducted (set in prog. result) 
		L641 long storageCost = GasCost.multiply(GasCost.CREATE_DATA, codeLength);
	* L881: has its own spendgas(gas, cause){.. 
		which in turn calls programResult.spendGas()

- ProgramResults.Java
	* gas is an argument.
	public void spendGas(long gas) {
        	gasUsed = GasCost.add(gasUsed, gas);
    	}

- VM.Java
	* L199 spendOpCodeGas(){
		if (!computeGas) { //default is true
	            return;
	        }
        	program.spendGas(gasCost, op.name());	
	}

	* Look at methods.. some of them like L303 doSMOD() just call spendOpCodeGas(), 
		but about 20 others e.g. L317 doEXP() compute gasCost using the "Gas Table".. 
	* computationally heave ones are charged, others like doADD are not.. (included in basic costs?)
	

