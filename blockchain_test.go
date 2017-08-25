package main 

import(
	"testing"
)

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
	b1 := &Block{Index: g.Index+1, PrevHash: g.Hash, Data: "Second", Hash: []byte{}}
	b1.calcHashForBlock()
	b2 := &Block{Index: g.Index, PrevHash: g.Hash, Data: "Second", Hash: []byte{}}
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

func TestIsValidChain(t *testing.T) {
	var chain Blockchain
	g := genesisBlock

	// base cases
	chain = Blockchain{[]Block{}}
	if chain.isValidChain(){
		t.Error("Validates empty chain")
	}

	chain = Blockchain{[]Block{g}}
	if chain.isValidChain(){
		t.Error("Validates chain of length 1")
	}

	// valid two block chain
	b1 := Block{Index: g.Index+1, PrevHash: g.Hash, Data: "Second", Hash: []byte{}}
	b1.calcHashForBlock()
	chain = Blockchain{[]Block{g, b1}}
	if !chain.isValidChain(){
		t.Error("does not accept valid chain")
	}

	//invalid two blockchain
	b2 := Block{Index: g.Index, PrevHash: g.Hash, Data: "Second", Hash: []byte{}}
	b2.calcHashForBlock()
	chain = Blockchain{[]Block{g, b2}}
	if chain.isValidChain(){
		t.Error("validates invalid chain of length two")
	}

	//valid three block chain
	b3 := Block{Index: b1.Index+1, PrevHash: b1.Hash, Data: "Third", Hash: []byte{}}
	b3.calcHashForBlock()
	chain = Blockchain{[]Block{g, b1, b3}}
	if !chain.isValidChain(){
		t.Error("fails to validate valid chain of length 3")
	}

	// invalid three block chain
	b4 := Block{Index: b2.Index+1, PrevHash: b2.Hash, Data: "Third", Hash: []byte{}}
	b4.calcHashForBlock()
	chain = Blockchain{[]Block{g, b2, b4}}
	if chain.isValidChain(){
		t.Error("validates invalid chain of length three")
	}
}




























