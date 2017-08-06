package main

import (
	"crypto/sha256"
	"encoding/binary"
)

type Block struct {
	Index uint32
	PrevHash []byte
	Info string
	Hash []byte
}

var genesisBlock = Block{0, []byte{0}, "genesis", []byte{0}}

// computes the hash of a block 
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

// used for comparison of hash byte slices
func testEqByteSlice (a, b []byte) bool {
    if a == nil && b == nil { 
        return true; 
    }
    if a == nil || b == nil { 
        return false; 
    }
    if len(a) != len(b) {
        return false
    }
    for i := range a {
        if a[i] != b[i] {
            return false
        }
    }
    return true
}













