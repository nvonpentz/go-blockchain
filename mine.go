package main 

import "fmt"

func listenToMinedBlockChannel(minedBlockChannel chan Block, blockWrapperChannel chan *BlockWrapper, myNode *Node,) {
	for {
        block := <- minedBlockChannel // user entered some input
        handleMinedBlock(block, minedBlockChannel, blockWrapperChannel, myNode)
	}
}

func handleMinedBlock(block Block, minedBlockChannel chan Block, blockWrapperChannel chan *BlockWrapper, n *Node) {
    if n.blockchain.isValidBlock(block){
        n.blockchain.addBlock(block)
        n.seenBlocks[string(block.Hash)] = true // specify weve now seen this block but don't update the blockWrapper address until its processed there
        go n.sendBlockWrapperFromMinedBlock(block, blockWrapperChannel)
    } else {
        fmt.Printf("Did not add mined block #%v\n", block.Index)
    }
    go n.blockchain.mineBlock(minedBlockChannel)
}