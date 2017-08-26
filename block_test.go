package main

import(
	"testing"
	// "fmt"
)

func areEqualBlocks(b1 Block, b2 Block) bool {
	indexEq    := b1.Index == b2.Index
	prevHashEq := string(b1.PrevHash) == string(b2.PrevHash)
	DataEq     := string(hashPacketList(b1.Data)) == string(hashPacketList(b2.Data))
	hashEq     := string(b1.Hash) == string(b2.Hash)

	return indexEq && prevHashEq && DataEq && hashEq
}

func TestIsValidNextBlock(t *testing.T){
	difficulty = 4294967295 // all hashses pass

	// test two equal blocks
	g  := &genesisBlock
	// b0 := &Block{}
	b1 := &Block{Index: g.Index + 1,
				 Nonce: 5000,
				 PrevHash: g.Hash,
				 Data: []Packet{},
				 Hash: []byte{}}

	b1.Hash = b1.calcHashForBlock(5000)

	// test valid block
	if !g.isValidNextBlock(b1){
		t.Error("Fails to validate valid next block")
	}

	// //test
	// if b0.isValidNextBlock(b0){
	// 	t.Error("Validates illegal (equal) blocks")
	// }

	// if g.isValidNextBlock(b0){
	// 	t.Error("Validates illegal blocks")
	// }

}


// func TestIsValidNextBlock(t *testing.T){
// 	// test two equal block
// 	empty := emptyBlock()
// 	b0 := &empty
// 	if b0.isValidNextBlock(b0){
// 		t.Error("Same blocks are validated as next block")
// 	}

// 	// test valid block
// 	g  := &genesisBlock
// 	b1 := &Block{Index: g.Index+1, PrevHash: g.Hash, Data: "Second", Hash: []byte{}}
// 	b1.calcHashForBlock()
// 	if !g.isValidNextBlock(b1){
// 		t.Error("Fails to validate valid block")
// 	}

// 	// test block whose index is wrong
// 	b2 := &Block{Index: g.Index, PrevHash: g.Hash, Data: "Second", Hash: []byte{}}
// 	b2.calcHashForBlock()
// 	if g.isValidNextBlock(b2){
// 		t.Error("Validates block whose index is incorrect")
// 	}
	
// 	// test block whose prevhash is wrong
// 	b3 := &Block{Index: g.Index+1, PrevHash: b2.Hash, Data: "Second", Hash: []byte{}}
// 	b3.calcHashForBlock()
// 	if g.isValidNextBlock(b3){
// 		t.Error("Validates block whose index is incorrect")
// 	}
// 	//test block whos hash is wrong (to be included when POW is added)
// 	// b3 := &Block{Index: g.Index + 1, PrevHash: g.Hash, Data: "Second", Hash: []byte{}}
// 	// if g.isValidNextBlock(b3){
// 	// 	t.Error("Validates block whose hash is incorrect")
// 	// }
// }

