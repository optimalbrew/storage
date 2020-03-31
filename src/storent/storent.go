package main

import(
	"time"
	"math"
	"fmt"
	//"github.com/ethereum/go-ethereum/rlp"
	//"bytes"
)

//storage cost per byte per sec
var rentRate = 1/math.Pow(2,21) //cannnot declare constant using a function. Value needs to be known at compile time. 

const(
	rentTrigIfMod = 1000 //collect if cell modified and rent due > this value
	rentTrigNotMod = 10000 //collect if cell NOT modified but rent due > this value

	sixMonths = 6*30*24*3600 // use for 6 months advance for new nodes  
)

//https://github.com/rsksmart/rskj/blob/master/rskj-core/src/main/java/co/rsk/trie/Trie.java
// the int sizes and types (uint) can be revised to minimize size
type trieNode struct{
	//... lots of things in RSKJ starting with arity = 2  //
	value 	[]byte
	left, right *trieNode //left and right children 
						  // left and right NodeReferences (separate class in the package on rskj)
	hash 	[]byte
	valueLength int64 //Uint24 in RSK // i.e. nodeSize
					  //if valueLength > 32 (bytes) && value==nil, then value has not been retrieved from state.
					  //code longer than 32 bytes is spread out over multiple nodes
	valueHash []byte
	childrenSize int64 //size of this node AND all its children! RSK VarInt to permit this to exceed 4GB
	//store of type TrieStore to store and retrieve nodes from trie

	//Make mods here
	leafnode bool // is leaf node (may not be necessary if RSK uses leaf prefix like Eth does 2,3)
	rentOutStanding int64 // can be negative, when some rent is pre-paid
	timeRentLastUpdated time.Time // alternative is to use Unix time (int64?)
}


// rent due: nodesize * time * rate (add overhead of 128 bytes to node size)
func rentBase(nodeSize, duration int64) int64 {
	return int64(rentRate * float64((nodeSize + 128) * duration)) //gas units
}


func main(){
	//timeDelta := time.Now().Unix() - 1585586197 //some arbitrary time
	//fmt.Printf("The difference is %v\n", timeDelta)
	//fmt.Printf("Storage gas cost is %.6f\n per byte per second", rentRate)
	fmt.Printf("Storage gas cost for 10 bytes for 6 months %d gas units.\n", rentBase(10, sixMonths))
	var test []byte = []byte{1,2}
	exNode := trieNode{value: test, valueLength: 10, leafnode: true}
	fmt.Printf("%+v\n", exNode)
	
}
