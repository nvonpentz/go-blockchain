package main 

import(
	// "fmt"
	"testing"
)

func generateMockChain() Blockchain {
	// create mock blockchain for use
	difficulty = 4294967295 // all hashses pass

	// create 4 different valid packets
	keys01 := GenerateNewKeypair()
	keys02 := GenerateNewKeypair()
	keys03 := GenerateNewKeypair()
	keys04 := GenerateNewKeypair()

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
	keys     := GenerateNewKeypair()
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
	chain := generateMockChain()
	lastBlock := chain.Blocks[4]

	if string(chain.getLastBlock().Hash) != string(lastBlock.Hash){
		t.Error("Fails to return last block in chain")
	}
}

func TestAddBlock(t *testing.T){
	chain := generateMockChain()
	lastBlock := chain.getLastBlock()

	b5 := &Block{Index: lastBlock.Index + 1,
				 Nonce: 5000,
				 PrevHash: lastBlock.Hash,
				 Data: []Packet{},
				 Hash: []byte{}}
	b5.Hash = b5.calcHashForBlock(5000)

	chain.addBlock(*b5)

	newLastBlock := chain.getLastBlock()

	if string(newLastBlock.Hash) != string(b5.Hash){
		t.Error("Fail to add new block")
	}
}

func TestIsValidChain(t *testing.T){
	chain := generateMockChain()

	// test valid chain
	if !chain.isValidChain(){
		t.Error("fails to validate valid chain")
	}

	g := chain.Blocks[0]

	// add in block to make invalid
	b5 := &Block{Index: g.Index + 1,
				 Nonce: 5000,
				 PrevHash: g.Hash,
				 Data: []Packet{},
				 Hash: []byte{}}
	b5.Hash = b5.calcHashForBlock(5000)

	invalidChain := chain
	invalidChain.Blocks = append(invalidChain.Blocks, *b5)

	if invalidChain.isValidChain(){
		t.Error("validates invalid chain")
	}
}







