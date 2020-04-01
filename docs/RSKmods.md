## Modifications to RSKIP113

**Motivation:** At present users only pay for the *amount of data* stored in blockchain state but *not the time* for which state values are stored. Storage rent will provide users with price signals to use state storage more efficiently.

### Summary
**Who pays:** storage rent is paid per node by their respective owners.

**How much:** $1/2^{21}$ gas units per byte per second.

**Who is it paid to:** It is collected by miners. Storage rent is an *additional, uncapped* source of revene for miners. Unlike transaction fees, there is no block gas limit for storage.  


### A. Introduce two new fields into Trie nodes
1. `timeRentLastUpdated`: Most recent time when storage rent was last computed for this node. 
- This can be used to measure time elasped between successive rent computations.
- This is independent of whether any rent payment was collected at that time. 
2. `rentOutStanding`: How much rent was outstanding at the last update (after accounting for any payments made at that time). 
- including any past amount 
- this can be negative, which indicates that some rent is pre-paid (paid in advance). In fact, it can even be encouraged (combined with an appropriate scheme to handle refunds).
- this can have a *upper bound* like 10,000 gas units (as in RSKIP113) 
- This is a *soft* upper bound because the outstanding rent can accumulate and exceed this value between node updates. 
- However, any time this node is touched by a transaction, it cannot be updated without a payment which leaves the outstanding rent under the threshold.
- This can help users estimate a time to live (TTL) every time they make a change to storage.

*Account hibernation:* A different threshold can be adopted (later) for account hibernation. The threshold can be used in combination with a lower bound on time duration since last update. 


### Proposed pseudocode

Let `dest` be any node that a TX will *modify* or *create* (reading is always free). This includes nodes that the TX sender does not own.


```
10. const storageRent = 1/(2^21) //gas per byte-sec

20. const rentCollectTrigger = 10000  //soft uppper bound on past due rent
30. const minPay = 1000               //avoid small transactions

// past rent due, whenever last updated (can be negative, prepaid)
40. RentPast = dest.rentOutStanding

//storage accumulated since last update
50. RentAccumulated = (time.Now() - dest.timeRentLastUpdated) * dest.prev_nodesize * storageRent

// optional TTL = desired time to live, advance payment. Min 6 months for new nodes
60. RentFuture = TTL * dest.new_nodeSize * storageRent

// add all sources
70. useRent = RentPast + RentAccum + RentFuture 

// minimum rent to pay now at end of TX
80. if (useRent - rentCollectTrigger > 0){  // e.g. rent due > 10k gas 
90.    minRent = max(useRent - rentCollectTrigger, minPay)  
100. } elseif (useRent > minPay) {                           
110.    minRent = minPay            // pay the mininum amount allowed 
120. } else {
130.        minRent = 0                 // pay nothing
    }

// deduct rent from sender provided it (rent offered) exceeds minimum
140. if (rentOfferred >= minRent){
150.    if (rentOfferred >= minPay){
160.        consumeRent(rentOffered)
170.        dest.rentOutStanding = useRent - consumedRent
            dest.timeRentLastUpdated = time.Now()
180.    } else{     // rent too low to consume (e.g. below 1000 gas)
            dest.rentOutStanding = useRent
            dest.timeRentLastUpdated = time.Now()
        }       
    } else {
        revertTransaction() // revert for all nodes referenced in the transaction
    }


```


### Concerns with previous proposal and pseudocode
1. RSKIP113 proposes a new field `lastRentPaidTime`. However, this -- by itself -- cannot completely account for accumulated storage rent.
- this only tracks payments, but does not track changes in nodesize (storage amount) between rent payments.
- example: a new node of size 10 bytes is created and is charged 6 months rent in advance. The `lastRentPaidTime` is set 6 months away. This means changes in the nodes size are not taken into account for rent caluclations for that time period.  



2. Rent to be collected even if node (`dest`) was not modifed. This is to be done only if rent exceeds 10k. 
- If the node is not modified, users can find some other way to use the information without paying rent (reading is always free).

### Other concerns
1. A token transfer changes storage (values, if not size) of atleast two accounts. If these transfers do not affect nodesize, no problem. Suppose, that a TX changes the receipient's storage size, how will that impact the rent to be collected for the recipient?
2. With the additional field `rentOutStanding`, we can force a computation/update. The owner of the account will then be charged in the future. Even the current sender need not be charged. What is key is that we must update  `rentOutStanding` and `timeRentLastUpdated`.
3. New accounts: is 6 months too much? What about users who use HD wallets to generate a sequence of single use accounts? Some rent must be charged, else users can split storage across accounts and or time to reduce payments. 

## Refunds
The refunds associated with `SSTORE` operations for deleting storage will need to account for any oustanding storage rent (including any rent paid in advance).

