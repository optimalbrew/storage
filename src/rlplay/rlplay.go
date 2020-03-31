package rlplay
//Example to illustrate use of RLP. Encode and devcode hex bytes string

import (
	"github.com/ethereum/go-ethereum/rlp"
	"encoding/hex"
	"fmt"
	"bytes"
)

// example of a struct for RLP encode/decode
type myStruct struct {
	A, B   uint
	String string
}

//encode struct to rlp Hex encoding
func EncodeMyStruct(m *myStruct) (string, error) {
	encoded, err := rlp.EncodeToBytes(&m)
	if err != nil{
		fmt.Printf("Error: %v\n", err)
		return "", err
	} else{
		return hex.EncodeToString(encoded), err
	}
}

//decode from RLP hex to struct
func DecodeMyStruct(rlpHex string)(myStruct, error){
	input, _ := hex.DecodeString(rlpHex)

	var s myStruct
	err := rlp.Decode(bytes.NewReader(input), &s)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return s, err
	} else {
		return s, err
	}

}

// moved most of this into tests (which will provide a main, change package name to main)
/* func main(){
	fmt.Println(Hello())

	//RLP decoding
	rlpHex := "d0150e8d6c75636b79206e756d62657273" 
	decoded, err := DecodeMyStruct(rlpHex)

	if err != nil {
    	fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Decoded value: %+v\n", decoded) //%+v : print struct field names
	} 

	//RLP encoding
	s_new := myStruct{A:21, B:14, String: "lucky numbers"}

	encodedStr, errStr := EncodeMyStruct(&s_new)
	if errStr != nil{
		fmt.Printf("Error: %v\n", errStr)
	} else{
		fmt.Printf("Encoded value: %v\n", encodedStr)
	}
	
} */
