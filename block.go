package main

import ("fmt"
		"time"
		// "math/rand"
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
	if blockchain.isValidBlock(block) == true {
		blockchain.Blocks = append(blockchain.Blocks, block)
		fmt.Printf("Added block %d to blockchain \n", block.Index)
	} else {
		fmt.Println("Did not add block to blockchain")
	}
}

func (blockchain *Blockchain) isValidBlock(block Block) bool{
	lastBlockInChain := blockchain.getLastBlock()
	if  lastBlockInChain.Index + 1 != block.Index{
		return false
	} else {
		return true
	}
}

func (blockchain *Blockchain) isValidChain() bool {
	blockchainLength := len(blockchain.Blocks)
	for i:= blockchainLength-1; i<=1; i-- {
		b2 := blockchain.Blocks[i]
		b1 := blockchain.Blocks[i-1]
		if isValidBlocks(b1, b2) == false {
			return false
		}
	}
	return true
}

func (blockchain Blockchain) getLastBlock() Block{
	lastBlock := blockchain.Blocks[len(blockchain.Blocks) - 1]
	return lastBlock
}

func (blockchain Blockchain) mineBlock(blockChannel chan Block, transmissionChannel chan *Transmission){
	fmt.Println("..begin mining..")
	// sleepTime := time.Duration((rand.Int() % 10) + 5) //use randomness for now
    time.Sleep(time.Second * 2)
	newBlockIndex := blockchain.getLastBlock().Index + 1
	newBlock := Block{newBlockIndex,"new block!"}
	fmt.Printf("Mined block # %d ", newBlock.Index)
    // trans := Transmission{newBlock, map[string]bool{}}
    // transmissionChannel <- &trans
	blockChannel <- newBlock
}

func isValidBlocks(b1 Block, b2 Block) (bool){
	if b2.Index != b1.Index + 1 {
		return false
	} else {
		return true
	}
}


