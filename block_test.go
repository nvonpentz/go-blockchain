package main

import(
	"testing"
)

func areEqualBlocks(b1 Block, b2 Block) bool {
	indexEq    := b1.Index == b2.Index
	prevHashEq := byteSlicesEqual(b1.PrevHash, b2.PrevHash)
	DataEq     := b1.Data == b2.Data
	hashEq     := byteSlicesEqual(b1.Hash, b2.Hash)

	return indexEq && prevHashEq && DataEq && hashEq
}

func TestIsValidNextBlock(t *testing.T){
	// test two equal block
	empty := emptyBlock()
	b0 := &empty
	if b0.isValidNextBlock(b0){
		t.Error("Same blocks are validated as next block")
	}

	// test valid block
	g  := &genesisBlock
	b1 := &Block{Index: g.Index+1, PrevHash: g.Hash, Data: "Second", Hash: []byte{}}
	b1.calcHashForBlock()
	if !g.isValidNextBlock(b1){
		t.Error("Fails to validate valid block")
	}

	// test block whose index is wrong
	b2 := &Block{Index: g.Index, PrevHash: g.Hash, Data: "Second", Hash: []byte{}}
	b2.calcHashForBlock()
	if g.isValidNextBlock(b2){
		t.Error("Validates block whose index is incorrect")
	}
	
	// test block whose prevhash is wrong
	b3 := &Block{Index: g.Index+1, PrevHash: b2.Hash, Data: "Second", Hash: []byte{}}
	b3.calcHashForBlock()
	if g.isValidNextBlock(b3){
		t.Error("Validates block whose index is incorrect")
	}
	//test block whos hash is wrong (to be included when POW is added)
	// b3 := &Block{Index: g.Index + 1, PrevHash: g.Hash, Data: "Second", Hash: []byte{}}
	// if g.isValidNextBlock(b3){
	// 	t.Error("Validates block whose hash is incorrect")
	// }
}

