package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

type Block struct {
	Index     uint32
	Nonce	  uint32
	PrevHash []byte
	Data     []Packet
	Hash     []byte
}

/* When sending a block to the main channel,
   we keep track of the sender in order to make requests
   for the entire blockchain if necessary */
type BlockWrapper struct {
    Block Block
    Sender string
}

var genesisBlock = Block{Index: 0, PrevHash: []byte{0}, Data: []Packet{}, Hash: []byte{0}}

func (block *Block) calcHashForBlock(nonce uint32) []byte {
	h := sha256.New()

	// convert nonce to bytes
	nonceBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(nonceBytes, nonce)

	// convert block index to bytes
	blockIndex := make([]byte, 4)
	binary.LittleEndian.PutUint32(blockIndex, block.Index)

	// hash the block data
	blockPacketsHash := hashPacketList(block.Data)

	h.Write(blockIndex)
	h.Write(block.PrevHash)
	h.Write(blockPacketsHash)
	h.Write(nonceBytes)
	
	return h.Sum(nil)
}

func (oldBlock *Block) isValidNextBlock(newBlock *Block) (bool){
	// new block's index must be one greater
	isValidIndex := newBlock.Index == oldBlock.Index + 1
	// fmt.Printf("isValidIndex %v \n", isValidIndex)

	// new block's previous hash has to equal the hash of the old block
	isValidPrevHash := string(newBlock.PrevHash) == string(oldBlock.Hash)
	// fmt.Printf("isValidPrevHash %v \n", isValidPrevHash)

	// all packets in block data must be valid
	areValidPacketSignatures := verifyPacketList(newBlock.Data)
	// fmt.Printf("areValidPacketSignatures %v \n", areValidPacketSignatures)

	// hash value must be below difficulty
	var newBlockHashAsInt uint32
	if len(newBlock.Hash) == 0 {
		fmt.Println(len(newBlock.Hash))
		fmt.Println("No hash, block invalid")
		return false
	} else {
		newBlockHashAsInt = binary.LittleEndian.Uint32(newBlock.Hash)		
	}

	isHashBelowDifficulty := newBlockHashAsInt < difficulty
	// fmt.Printf("newBlockHashAsInt %v, difficulty %v", newBlockHashAsInt, difficulty)
	// fmt.Printf("isHashBelowDifficulty %v \n", isHashBelowDifficulty)


	// hash of entire block must equal the claimed block hash
	calculatedBlockHash := newBlock.calcHashForBlock(newBlock.Nonce)
	isCorrectBlockHash  := string(calculatedBlockHash) == string(newBlock.Hash)
	// fmt.Printf("isCorrectBlockHash %v \n", isCorrectBlockHash)

	isValidBlock := isValidIndex &&
					isValidPrevHash &&
					areValidPacketSignatures &&
					isHashBelowDifficulty &&
					isCorrectBlockHash

	//this is where proof of work comes to validate the calculated hash
	return isValidBlock
}













