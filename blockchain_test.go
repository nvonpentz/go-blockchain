package main 

import(
	"testing"
	"fmt"
)

func TestIsValidBlock(t *testing.T){
	b1 := Block{0, []byte{}, "", []byte{}}
	b2 := Block{0, []byte{}, "", []byte{}}

	chain := Blockchain{[]Block{b1}}
	if chain.isValidBlock(b2) != false {
		t.Error("Validates illegal block")
	}

	b1 = Block{0, []byte{}, "", []byte{}}
	b2 = Block{1, []byte{byte('3'), byte('4')}, "", []byte{}}
	chain = Blockchain{[]Block{b1}}
	if chain.isValidBlock(b2) != false {
		t.Error("Validates block that where hashes don't align")
	}

	b1 = Block{0, []byte{}, "", []byte{byte('3'), byte('4')}}
	b2 = Block{1, []byte{byte('3'), byte('4')}, "", []byte{}}
	chain = Blockchain{[]Block{b1}}
	if chain.isValidBlock(b2) != true {
		t.Error("Does not approve valid block")
	}
}

func TestGetLastBlock(t *testing.T){
	b1    := Block{0, []byte{}, "", []byte{byte('3'), byte('4')}}
	b2    := Block{1, []byte{byte('3'), byte('4')}, "", []byte{}}
	chain := Blockchain{[]Block{b1,b2}}

	lb    := chain.getLastBlock()

	fmt.Println(b2)
	fmt.Println(lb)

	if areEqualBlocks(b2, lb) != true {
		t.Error("Failed to get last block")
	}
}

