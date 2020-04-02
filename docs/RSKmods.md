## Modifications to rent proposal


### Concerns with pseudocode in RSKIP113

```
10. if (d>lastRentPaidTime) {
20.    useRentGas =  nodeSize*(d-lastRentPaidTime)/2^21
30A    if ((dest was modified) && (useRentGas>=1000)) || 
30B.       ((dest was NOT modified) && (useRentGas>=10000)) {
40.        dest.lastRentPaidTime = now
50.        consumeRentGas(useRentGas);
      }
}
```
This pseudocode in RSKIP 113 was meant to be illustrative, not exhaustive. Thus, it does not include some details present in the text of the proposal, e.g. advanced rent of 6 months to be collected from new nodes or deducting some gas even if a transaction is reverted. 

However, there are other concerns:
1. The computation of `useRentGas` (line 20) does not take into account that some rent may have been paid in advance or that some rent may be outstanding from before.

2. Not every change in a node results in a rent payment. A rent payment is triggerred only if condition in lines 30A,B is met. The rent computation **does not track changes** in `nodesize` in between two rent payments.  

3. Advance payments allow user to predict a minimum time to live (TTL) for the account. But the current logic  `consumeRentGas(useRentGas)` (line 50) does not permit advanced payments for existing nodes.

4. If a new node is created, then we enter the `else` block for `if` (line 10), which is not specified. 
- Suppose new node is charged 6 months rent in advance. The `lastRentPaidTime` is set 6 months away. This means changes in the `nodesize` may not be taken into account for rent calculations for 6 months.

*My takeaway:* By itself, `lastRentPaidTime` seems inadequate for rent calculations and payments. It may be adequate if nodesize was constant or if advance payments are not allowed.


### Suggestion 
Introduce two new fields into Trie nodes (instead of `lastRentPaidTime`)
1. `timeRentLastUpdated`: Most recent time when storage rent was last **computed** (not paid) for this node. 
- This can be used to measure time elasped between successive rent computations.
- This is independent of whether any rent payment was collected at that time. 
2. `rentOutStanding`: How much rent was outstanding at the last update (after accounting for any payments made at that time).  
- this can be negative, which indicates that some rent is pre-paid (paid in advance). In fact, it can even be encouraged (combined with an appropriate scheme to handle refunds).
- this can have a *upper bound* like 10,000 gas units (as in RSKIP113) 
- This is a *soft* upper bound because the outstanding rent can accumulate and exceed this value between node updates. 
- However, any time this node is modified by a transaction, it cannot be updated without a payment which leaves the outstanding rent under the threshold.
- This can help users control the time to live (TTL) each time they make a change to storage.

*Account hibernation:* Some combinations of these fields can be used (later) to set triggers for account hibernation. 


### Proposed pseudocode
This proposal is also meant to be illustrative, not exhaustive.


Let `dest` be any node that a currently executing transaction `TX_current` will *modify* or *create*. 

```
10. const storageRent = 1/(2^21) //gas per byte-sec
20. const rentCollectTrigger = 10000  // A "soft" limit on outstanding rent
30. const dontPayIfBelow = 1000       // min. payment to avoid small transactions

/* 3 possible components of rent
    1. past rent due, whenever last updated i.e. calculated (can be negative, if prepaid)
    2. storage rent accumulated since last update (use previous nodesize)
    3. Future rent paid in advance. TTL (time to live), min 6 months for new nodes
    Future rent depends on NEW nodesize 
*/
40. RentPast = dest.rentOutStanding
50. RentAccumulated = (time.Now() - dest.timeRentLastUpdated) * dest.previous_nodesize * storageRent
60. RentFuture = TTL * dest.updated_nodesize * storageRent     //TTL, nodesize from TX_current or 6 months (new node)

// add all sources
70. useRentGas = RentPast + RentAccumulated + RentFuture 

// minimum rent to pay now at end of TX
80. if (useRentGas - rentCollectTrigger > 0){   // rent exceeds Trigger (e.g. 10k) 
90.    minRentNow = max(useRentGas - rentCollectTrigger, dontPayIfBelow)  
100. } elseif (useRentGas > dontPayIfBelow) {                           
110.    minRentNow = dontPayIfBelow             // min allowed (e.g. 1k) 
120. } else {
130.    minRentNow = 0                          // pay nothing at present
    }

// deduct rent from sender provided rent offered via TX exceeds minimum
140. if (rentOfferred >= minRentNow){   // rentOfferred from TX_curent
150.    if (rentOfferred >= dontPayIfBelow){
160.        consumeRent(rentOffered)
170.        dest.rentOutStanding = useRentGas - consumedRent  //update trie
180.        dest.timeRentLastUpdated = time.Now()             // update trie
190.    } else{     // rent too low to consume (e.g. below 1k gas)
200.        dest.rentOutStanding = useRentGas       // update trie
210.        dest.timeRentLastUpdated = time.Now()   // update trie
        }       
    } else {
220.        revertTransaction() // revert changes/updates for ALL nodes in TX_current
    }

```

### Notes

20: This soft limit can never be exceeded as a result of some transaction. Only passive rent accumulation (no change in node for long periods) can push the outstanding rent beyond this value.


220: This includes nodes that the transaction sender does not own and **may not** be responsible to pay rent. However, we do not want to allow sender to impose arbitrary rent on other accounts. So we adopt a conservative approach. The outstanding rent for accounts cannot be pushed beyond the soft limit. **Issue:** How to handle impact on others? Should be left to contract designer, app developer.




## Other Issues
- The refunds associated with `SSTORE` operations for deleting storage will need to account for any outstanding storage rent (including any rent paid in advance).
- Rent is paid to miners. Storage costs are perhaps an insignificant component of overall expenses of running a mining operation. That is not the case for independent (non mining) full nodes. It may be beneficial for the ecosystem to develop ways to subsidize the costs of running indepdent full nodes.

