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

func (blockchain *Blockchain) verifyChain() bool {
	blockchainLength := len(blockchain.Blocks)
	for i:= blockchainLength-1; i<=1; i-- {
		b2 := blockchain.Blocks[i]
		b1 := blockchain.Blocks[i-1]
		if verifyBlocks(b1, b2) == false {
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
	sleepTime := time.Duration((rand.Int() % 10) + 5)
    time.Sleep(time.Second * sleepTime)
	newBlockIndex := blockchain.getLastBlock().Index + 1
	newBlock := Block{newBlockIndex,"new block!"}
	// blockchain.Blocks = append(blockchain.Blocks, newBlock)
	fmt.Println("Mined a new block!")
	blockChannel <- newBlock
}

func verifyBlocks(b1 Block, b2 Block) (bool){
	if b2.Index != b1.Index + 1 {
		return false
	} else {
		return true
	}
}


