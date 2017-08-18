package main 

// import "fmt"

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
	// if blockchainLength == 2 {
	// 	firstBlock  := blockchain.Blocks[0]
	// 	secondBlock := blockchain.Blocks[1]
	// 	return firstBlock.isValidNextBlock(&secondBlock)
	// } else {
	// 	fmt.Println("not length two")
	// }

	for i:=blockchainLength-1; i>=1; i-- {
		b2 := blockchain.Blocks[i]
		b1 := blockchain.Blocks[i-1]
		if b1.isValidNextBlock(&b2) == false {
			return false
		} else{
		}
	}
	// fmt.Println("about to return")
	return true
}

func (blockchain Blockchain) getLastBlock() Block{
	lastBlock := blockchain.Blocks[len(blockchain.Blocks) - 1]
	return lastBlock
}