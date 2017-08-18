package main 

import(
	"testing"
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
	lastBlock    := chain.getLastBlock()

	if areEqualBlocks(b2, lastBlock) != true {
		t.Error("Failed to get last block")
	}
}

func TestAddBlock(t *testing.T){
	g  := genesisBlock
	b1 := &Block{Index: g.Index+1, PrevHash: g.Hash, Info: "Second", Hash: []byte{}}
	b1.calcHashForBlock()
	b2 := &Block{Index: g.Index, PrevHash: g.Hash, Info: "Second", Hash: []byte{}}
	b2.calcHashForBlock()

	chain1 := &Blockchain{[]Block{g}}
	chain1.addBlock(*b1)
	chain2 := &Blockchain{[]Block{g}}
	chain2.addBlock(*b2)


	// make sure valid is added
	if len(chain1.Blocks) != 2 {
		t.Error("failed to add valid block")
	}

	// make sure invalid is not added
	if len(chain2.Blocks) != 1 {
		t.Error("added invalid block when we shouldn't")
	}

}

// func TestIsValidChain(t *testing.T) {
// 	chain1 := &Blockchain{[]Block{g}}
// }














