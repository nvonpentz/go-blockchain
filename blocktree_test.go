package main 

import(
	"testing"
)

func TestEqualBTNodes(t *testing.T){
	b1 := BTNode{0, &genesisNode, []byte{0}, "Data", []byte{0}}
	b2 := BTNode{0, &genesisNode, []byte{0}, "Data", []byte{0}}

	eq := equalBTNodes(b1, b2)

	if eq != true {
		t.Error("Identical nodes shoul be equal")
	}
}

func TestCalcBTNodeHash(t *testing.T) {}

func TestIsValidNextBTNode(t *testing.T){
	b1      := BTNode{0, &genesisNode, []byte{0}, "Data", []byte{0}}
	b1.Hash  = b1.calcBTNodeHash()
	b2valid := BTNode{1, &b1, b1.Hash, "Data", []byte{0}}

	if b1.isValidNextBTNode(&b2valid) != true {
	 	t.Error("Valid node does not validate")
	 }
}