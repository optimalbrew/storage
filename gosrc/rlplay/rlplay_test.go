package rlplay

import "testing"

func TestEncodeMyStruct(t *testing.T) {
    s_new := myStruct{A:21, B:14, String: "lucky numbers"}
	got, _ := EncodeMyStruct(&s_new)    
    want := "d0150e8d6c75636b79206e756d62657273"
    if got != want {
        t.Errorf("Encoded: %v, wanted %v", got, want)
    }
}

func TestDecodeMyStruct(t *testing.T) {
    rlpHex := "d0150e8d6c75636b79206e756d62657273"
	got, _ := DecodeMyStruct(rlpHex) 
    want := myStruct{A:21, B:14, String: "lucky numbers"}
    if got != want {
        t.Errorf("Decoded: %v, wanted %v", got, want)
    }
}