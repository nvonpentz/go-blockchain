package main 

import "fmt"

type Blockchain struct {
	Blocks []Block
}

func (b Blockchain) printBlockchain(){
    for i := range b.Blocks {
        block := b.Blocks[i]
        fmt.Printf("  Block %d is: \n   PrevHash: %v \n   Info:     %v \n   Hash:     %v \n", i, block.PrevHash, block.Info, block.Hash)
    }
}

func (blockchain *Blockchain) addBlock(block Block) {
	lastBlock := blockchain.getLastBlock()

	if lastBlock.isValidNextBlock(&block) == true {
		blockchain.Blocks = append(blockchain.Blocks, block)
	} else {
	}
}

func (blockchain *Blockchain) isValidChain() bool {
	blockchainLength := len(blockchain.Blocks)
	if blockchain == nil || blockchainLength == 0 || blockchainLength == 1 { return false }

	for i:=blockchainLength-1; i>=1; i-- {
		b2 := blockchain.Blocks[i]
		b1 := blockchain.Blocks[i-1]
		if b1.isValidNextBlock(&b2) == false {
			return false
		} else{
		}
	}
	return true
}

func (blockchain Blockchain) getLastBlock() Block{
	lastBlock := blockchain.Blocks[len(blockchain.Blocks) - 1]
	return lastBlock
}