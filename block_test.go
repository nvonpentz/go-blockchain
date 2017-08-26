package main

import(
	"testing"
	"github.com/nvonpentz/go-hashable-keys"
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

	// create 3 different valid packets
	keys01 := hashkeys.GenerateNewKeypair()
	keys02 := hashkeys.GenerateNewKeypair()
	keys03 := hashkeys.GenerateNewKeypair()


	packet01 := createPacket("document.txt", *keys01)
	packet02 := createPacket("document.txt", *keys02)
	packet03 := createPacket("document.txt", *keys03)

	packets := []Packet{packet01, packet02, packet03}

	// test two equal blocks
	g  := &genesisBlock
	// b0 := &Block{}
	b1 := &Block{Index: g.Index + 1,
				 Nonce: 5000,
				 PrevHash: g.Hash,
				 Data: packets,
				 Hash: []byte{}}

	b1.Hash = b1.calcHashForBlock(5000)

	// test valid block
	if !g.isValidNextBlock(b1){
		t.Error("Fails to validate valid next block")
	}

	// test block with wrong index
	b2      := *b1
	b2.Index = g.Index 
	b2.Hash  = b2.calcHashForBlock(5000)
	if g.isValidNextBlock(&b2){
		t.Error("Validates block with incorrect index")
	}

	// test block with wrong prevHash
	b3 := *b1
	b3.PrevHash = b2.Hash // wrong hash
	b3.Hash = b3.calcHashForBlock(5000)
	if g.isValidNextBlock(&b3){
		t.Error("Validates block with incorrect prevHash")
	}

	// test block with incorrect hash
	b4 := *b1
	b4.Hash = b2.Hash //wrong hash
	if g.isValidNextBlock(&b4){
		t.Error("Validates block with incorrect hash")
	}

	// test block with invalid data
	// create invalid packet
	doc       := readDocument("document.txt")
	hashedDoc := hashDocument(doc)
	signature := signHash(hashedDoc, *keys01) // sign with keys01
	packet04 := Packet{Hash: hashedDoc, Signature: signature, Owner: keys02.Public}

	packets = []Packet{packet01, packet04, packet02, packet03}
	
	b5 := *b1
	b5.Data = packets
	b5.Hash = b5.calcHashForBlock(5000)
	if g.isValidNextBlock(&b5){
		t.Error("Validates block with invalid packets")
	}

	// test block who's hash doesn't meet difficulty target
	difficulty = 1 //impossibel
	if g.isValidNextBlock(b1){
		t.Error("Validates block that doesn't meet difficulty requirement")
	}	
}







