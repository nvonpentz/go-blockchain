package main

import (
	"crypto/sha256"
	"encoding/binary"
)

type Block struct {
	Index     uint32
	PrevHash []byte
	Info     string // The data stored on the block
	Hash     []byte
}

/* When sending a block to the main channel,
   we keep track of the sender in order to make requests
   for the entire blockchain if necessary */
type BlockWrapper struct {
    Block Block
    Sender string
}

func emptyBlock() Block{
	return Block{Index: 0, PrevHash: []byte{}, Info: "", Hash: []byte{}}
}

func emptyBlockWrapper() BlockWrapper{
	return BlockWrapper{Block: emptyBlock(), Sender: "127.0.0.1:1999"}
}

var genesisBlock = Block{0, []byte{0}, "genesis", []byte{0}}

func (block *Block) calcHashForBlock() []byte {
	blockHash := sha256.New()

	nbIndexBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(nbIndexBytes, block.Index)
	pbHashBytes  := []byte(block.PrevHash)
	nbInfoBytes  := []byte(block.Info)
	toHash := append(nbIndexBytes, pbHashBytes...)
	toHash  = append(toHash, nbInfoBytes...)
	blockHash.Write(toHash)

	return blockHash.Sum(nil)
}

func (oldBlock *Block) isValidNextBlock(newBlock *Block) (bool){
	// new block's index must be one greater
	isValidIndex := newBlock.Index == oldBlock.Index + 1

	// new block's previous hash has to equal the hash of the old block
	isValidPrevHash := byteSlicesEqual(newBlock.PrevHash, oldBlock.Hash)
	isValidBlock := isValidIndex && isValidPrevHash

	//this is where proof of work comes to validate the calculated hash
	return isValidBlock
}













