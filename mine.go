package main 

import ("fmt")

func handleMinedBlock(block Block, blockWrapperChannel chan *BlockWrapper, n *Node) {
    if n.blockchain.isValidBlock(block){
        n.blockchain.addBlock(block)
        n.seenBlocks[string(block.Hash)] = true // specify weve now seen this block but don't update the blockWrapper address until its processed there
        go n.sendBlockWrapperFromMinedBlock(block, blockWrapperChannel)
    } else {
        fmt.Printf("Did not add mined block #%v\n", block.Index)
    }
    go n.blockchain.mineBlock(blockWrapperChannel, n)
}