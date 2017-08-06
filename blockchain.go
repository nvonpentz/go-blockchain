package main 

import (
	"fmt"
	"time"
	"math/rand"
)

type Blockchain struct {
	Blocks []Block
}

func (blockchain *Blockchain) addBlock(block Block) {
	if blockchain.isValidBlock(block) == true {
		blockchain.Blocks = append(blockchain.Blocks, block)
	} else {
	}
}

func (blockchain *Blockchain) isValidBlock(block Block) bool{
	lastBlockInChain := blockchain.getLastBlock()
	return lastBlockInChain.isValidNextBlock(&block)
}

func (blockchain *Blockchain) isValidChain() bool {
	blockchainLength := len(blockchain.Blocks)
	if blockchain == nil || blockchainLength == 0 || blockchainLength == 1 { return false }
	for i:= blockchainLength-1; i<=1; i-- {
		b2 := blockchain.Blocks[i]
		b1 := blockchain.Blocks[i-1]
		if b1.isValidNextBlock(&b2) == false {
			return false
		}
	}
	return true
}

func (blockchain Blockchain) getLastBlock() Block{
	lastBlock := blockchain.Blocks[len(blockchain.Blocks) - 1]
	return lastBlock
}

func (blockchain *Blockchain) mineBlock(blockChannel chan Block){
	fmt.Println("-> begin mining..")

	// sleep between 5 - 10 seconds before mining block to simulate a blockchain
	sleepTime := time.Duration((rand.Int() % 10) + 5)
    time.Sleep(time.Second * sleepTime)

    //create new block
    prevBlock     := blockchain.getLastBlock()
	newBlockIndex := prevBlock.Index + 1
	newBlockInfo  := "new block!"
	newBlock := Block{newBlockIndex, prevBlock.Hash, newBlockInfo, []byte{}}

	// must calculate the hash of this block
	newBlockHash := newBlock.calcHashForBlock()
	newBlock      = Block{newBlockIndex, prevBlock.Hash, newBlockInfo, newBlockHash}

	// send to control center to 
	blockChannel <- newBlock 
}