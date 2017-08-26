package main 

import(
	// "fmt"
	"testing"
	"github.com/nvonpentz/go-hashable-keys"
)

func generateMockChain() Blockchain {
	// create mock blockchain for use
	difficulty = 4294967295 // all hashses pass

	// create 4 different valid packets
	keys01 := hashkeys.GenerateNewKeypair()
	keys02 := hashkeys.GenerateNewKeypair()
	keys03 := hashkeys.GenerateNewKeypair()
	keys04 := hashkeys.GenerateNewKeypair()

	packet01 := createPacket("document.txt", *keys01)
	packet02 := createPacket("document.txt", *keys02)
	packet03 := createPacket("document.txt", *keys03)
	packet04 := createPacket("document.txt", *keys04)

	packets01  := []Packet{packet01}
	packets02  := []Packet{packet02, packet03}
	packets03  := []Packet{packet03}
	packets04  := []Packet{packet04}

	g  := &genesisBlock
	b1 := &Block{Index: g.Index + 1,
				 Nonce: 5000,
				 PrevHash: g.Hash,
				 Data: packets01,
				 Hash: []byte{}}
	b1.Hash = b1.calcHashForBlock(5000)

	b2 := &Block{Index: b1.Index + 1,
				 Nonce: 5000,
				 PrevHash: b1.Hash,
				 Data: packets02,
				 Hash: []byte{}}
	b2.Hash = b2.calcHashForBlock(5000)

	b3 := &Block{Index: b2.Index + 1,
				 Nonce: 5000,
				 PrevHash: b2.Hash,
				 Data: packets03,
				 Hash: []byte{}}
	b3.Hash = b3.calcHashForBlock(5000)

	b4 := &Block{Index: b3.Index + 1,
				 Nonce: 5000,
				 PrevHash: b3.Hash,
				 Data: packets04,
				 Hash: []byte{}}
	b4.Hash = b4.calcHashForBlock(5000)

	chain := Blockchain{Blocks: []Block{*g, *b1, *b2, *b3, *b4}}

	return chain
}


func TestFindPacketByHash(t *testing.T){
	chain    := generateMockChain()
	packet02 := chain.Blocks[2].Data[0]
	// fmt.Printf("searching for packet hash: %v \n", packet02)

	p2       := chain.findPacketByHashAndPublicKey(packet02.Hash, packet02.Owner)
	if !equalPackets(p2, packet02){
		t.Error("Fails to find packet in blockchain via hash")
	}

	// find first packet in blockchain
	packet01 := chain.Blocks[1].Data[0]
	p1       := chain.findPacketByHashAndPublicKey(packet01.Hash, packet01.Owner)

	// generate new packet not in blockcahin
	keys     := hashkeys.GenerateNewKeypair()
	packet05 := createPacket("document.txt", *keys)
	if !equalPackets(p1, packet01){
		t.Error("Fails to find first packet in blockchain via hash")
	}

	p5 := chain.findPacketByHashAndPublicKey(packet05.Hash, packet05.Owner)
	if !equalPackets(p5, Packet{}){
		t.Error("returns non empty packet when searching for a non existing packet in blockchain")
	}
}

func TestGetLastBlock(t *testing.T){
	// TO DO
}

func TestAddBlock(t *testing.T){
	// TO DO
}

func TestIsValidChain(t *testing.T){
	// TO DO
}

// func TestAddBlock(t *testing.T){
// 	g  := genesisBlock
// 	b1 := &Block{Index: g.Index+1, PrevHash: g.Hash, Data: "Second", Hash: []byte{}}
// 	b1.calcHashForBlock()
// 	b2 := &Block{Index: g.Index, PrevHash: g.Hash, Data: "Second", Hash: []byte{}}
// 	b2.calcHashForBlock()

// 	chain1 := &Blockchain{[]Block{g}}
// 	chain1.addBlock(*b1)
// 	chain2 := &Blockchain{[]Block{g}}
// 	chain2.addBlock(*b2)


// 	// make sure valid is added
// 	if len(chain1.Blocks) != 2 {
// 		t.Error("failed to add valid block")
// 	}

// 	// make sure invalid is not added
// 	if len(chain2.Blocks) != 1 {
// 		t.Error("added invalid block when we shouldn't")
// 	}

// }

// func TestIsValidChain(t *testing.T) {
// 	var chain Blockchain
// 	g := genesisBlock

// 	// base cases
// 	chain = Blockchain{[]Block{}}
// 	if chain.isValidChain(){
// 		t.Error("Validates empty chain")
// 	}

// 	chain = Blockchain{[]Block{g}}
// 	if chain.isValidChain(){
// 		t.Error("Validates chain of length 1")
// 	}

// 	// valid two block chain
// 	b1 := Block{Index: g.Index+1, PrevHash: g.Hash, Data: "Second", Hash: []byte{}}
// 	b1.calcHashForBlock()
// 	chain = Blockchain{[]Block{g, b1}}
// 	if !chain.isValidChain(){
// 		t.Error("does not accept valid chain")
// 	}

// 	//invalid two blockchain
// 	b2 := Block{Index: g.Index, PrevHash: g.Hash, Data: "Second", Hash: []byte{}}
// 	b2.calcHashForBlock()
// 	chain = Blockchain{[]Block{g, b2}}
// 	if chain.isValidChain(){
// 		t.Error("validates invalid chain of length two")
// 	}

// 	//valid three block chain
// 	b3 := Block{Index: b1.Index+1, PrevHash: b1.Hash, Data: "Third", Hash: []byte{}}
// 	b3.calcHashForBlock()
// 	chain = Blockchain{[]Block{g, b1, b3}}
// 	if !chain.isValidChain(){
// 		t.Error("fails to validate valid chain of length 3")
// 	}

// 	// invalid three block chain
// 	b4 := Block{Index: b2.Index+1, PrevHash: b2.Hash, Data: "Third", Hash: []byte{}}
// 	b4.calcHashForBlock()
// 	chain = Blockchain{[]Block{g, b2, b4}}
// 	if chain.isValidChain(){
// 		t.Error("validates invalid chain of length three")
// 	}
// }















