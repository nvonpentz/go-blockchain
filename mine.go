package main 

import (
	"fmt"
	// "time"
	// "math/rand"
)

func mineBlock(blockWrapperChannel chan *BlockWrapper, n *Node){
	fmt.Println("-> begin mining")

	
	// // sleep between 5 - 10 seconds before mining block to simulate a blockchain
	// sleepTime := time.Duration((rand.Int() % 4) + 2)
 //    time.Sleep(time.Second * sleepTime)

 //    //create new block
 //    prevBlock     := n.blockchain.getLastBlock()
	// newBlockIndex := prevBlock.Index + 1
 //    fmt.Printf("newBlockIndex: %v\n:", newBlockIndex)
	// newBlockData  := "new block!"
	// newBlock := Block{newBlockIndex, prevBlock.Hash, newBlockData, []byte{}}

	// // must calculate the hash of this block
	// newBlockHash := newBlock.calcHashForBlock()
	// newBlock      = Block{newBlockIndex, prevBlock.Hash, newBlockData, newBlockHash}
 //    fmt.Printf("Mined block: %v\n", newBlock.Index)
 //    blockWrapperChannel <- &BlockWrapper{Block: newBlock, Sender: n.address}
	// handleMinedBlock(newBlock, blockWrapperChannel, n)
    mineBlock(blockWrapperChannel, n)
}






