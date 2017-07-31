package main

import ("fmt"
		"time"
		"math/rand"
		)

type Blockchain struct {
	Blocks []Block
}

type Block struct {
	Index int
	Information string
}

var genesisBlock = Block{0, "genesis transaction"}

func (blockchain *Blockchain) addBlock(block Block) {
	if blockchain.verifyBlock(block) == true {
		blockchain.Blocks = append(blockchain.Blocks, block)
		fmt.Println("Added block to blockchain")
	} else {
		fmt.Println("Did not add block to blockchain")
	}
}

func (blockchain *Blockchain) verifyBlock(block Block) bool{
	lastBlockInChain := blockchain.getLastBlock()
	if  lastBlockInChain.Index + 1 != block.Index{
		fmt.Println("Blockchain does not verify")
		return false
	} else {
		fmt.Println("Blockcahin verifies")
		return true
	}
}

func (blockchain Blockchain) getLastBlock() Block{
	lastBlock := blockchain.Blocks[len(blockchain.Blocks) - 1]
	return lastBlock
}

func (blockchain *Blockchain) mineBlock(blockChannel chan Block){
	sleepTime := time.Duration((rand.Int() % 10) + 5)
    time.Sleep(time.Second * sleepTime)
	newBlockIndex := blockchain.getLastBlock().Index + 1
	newBlock := Block{newBlockIndex,"new block!"}
	// blockchain.Blocks = append(blockchain.Blocks, newBlock)
	fmt.Println("Mined a new block!")
	blockChannel <- newBlock
}



