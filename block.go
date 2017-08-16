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

var genesisBlock = Block{0, []byte{0}, "genesis", []byte{0}}

// for testing
// func emptyBlock() Block{
// 	return Block{0, []byte{}, "", []byte{}}
// }

// for testing
func areEqualBlocks(b1 Block, b2 Block) bool {
	indexEq    := b1.Index == b2.Index
	prevHashEq := testEqByteSlice(b1.PrevHash, b2.PrevHash)
	infoEq     := b1.Info == b2.Info
	hashEq     := testEqByteSlice(b1.Hash, b2.Hash)

	return indexEq && prevHashEq && infoEq && hashEq
}

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
	isValidHash := testEqByteSlice(newBlock.PrevHash, oldBlock.Hash)
	isValidBlock := isValidIndex && isValidHash

	return isValidBlock
}













