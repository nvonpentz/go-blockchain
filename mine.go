package main 

import (
	"fmt"
	"time"
	"math/rand"
)

func handleMinedBlock(block Block, blockWrapperChannel chan *BlockWrapper, n *Node) {
    if n.blockchain.isValidBlock(block){
        n.blockchain.addBlock(block)
        n.seenBlocks[string(block.Hash)] = true // specify weve now seen this block but don't update the blockWrapper address until its processed there
        go n.sendBlockWrapperFromMinedBlock(block, blockWrapperChannel)
    } else {
        fmt.Printf("Did not add mined block #%v\n", block.Index)
    }
    go mineBlock(&n.blockchain, blockWrapperChannel, n)
}

func mineBlock(blockchain *Blockchain, blockWrapperChannel chan *BlockWrapper, n *Node){
	fmt.Println("-> begin mining..")

	// sleep between 5 - 10 seconds before mining block to simulate a blockchain
	sleepTime := time.Duration((rand.Int() % 3) + 1)
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
	handleMinedBlock(newBlock, blockWrapperChannel, n)

	// minedBlockChannel <- newBlock 
}