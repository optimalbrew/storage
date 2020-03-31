## Modifications based on RSKIP113

**Motivation:** At present users only pay for the *amount of data* stored in blockchain state but *not the time* for which state values are stored. Storage rent will provide users with price signals to use state storage more efficiently.

### Summary
**Who pays:** storage rent is paid per node. The sender of the transaction affecting the change (or creating a new node) pays the rent.
**How much:** $1/2^{21}$ gas units per byte per second.
- Note, RSK gas prices are in BTC not ETH.
**Who is it paid to:** It is collected by miners as part of block rewards, just like transaction fees. 
- Storage rent is an additional source of revene for miners. Unlike transaction fees, there is no block gas limit for storage.  


### A. Introduce two new fields into Trie nodes
1. `timeRentLastUpdated`: Most recent time when storage rent was last computed for this node. 
- This is independent of whether any rent payment was collected at that time.
- This can be used to measure time deltas. 
2. `rentOutStanding`: How much rent was outstanding after the last update (not last payment) 
- including any past amount 
- this can be negative, which indicates that some rent is pre-paid (paid in advance)
- this can have a *upper bound* like 10,000 gas units (as in RSKIP113) 
- This is a *soft* upper bound because the outstanding rent can accumulate and exceed this value between node updates. 
- However, any time this node is touched by a transaction, it cannot be updated without a payment which leaves the outstanding rent under the threshold. 
- *Hibernation:* A different threshold can be adopted (later) for account hibernation. The threshold can be used in combination with a lower bound on time duration since last update. 

### Why
RSKIP113 proposes a new field `lastRentPaidTime`. However, this -- by itself -- cannot completely account for accumulated storage rent.
- 


RSKIP pseudicode: collect rent if 
- the time since last payment  (was not doing anything, except splitting new nodes from old ones)

Threshold when node is not modified is 1k. 


### Pseudocode

Let `dest` be any node which the sender $S$ of a TX $T$ is trying to *modify* or *create* (reading is always free). 

```






```

