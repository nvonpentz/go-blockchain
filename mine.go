package main 

import (
	"fmt"
	"encoding/binary"
)

const difficulty = 20000

func mineBlock(blockWrapperChannel chan *BlockWrapper, n *Node){
	fmt.Println("-> begin mining...")

	var blockHashAsInt uint32
	var nonce          uint32
	nonce          = 0
	blockHashAsInt = 4294967295 // max value ensures we will enter mining loop

	var lastBlock      Block
	var block          Block
	var currentPackets []Packet
	var blockHash      []byte

	for blockHashAsInt > difficulty {
		lastBlock 	    = n.blockchain.getLastBlock()
		currentPackets  = n.curPacketList
		block 		    = Block{Index:    lastBlock.Index + 1,
								Nonce:    nonce,
								PrevHash: lastBlock.Hash,
								Data:     currentPackets,
								Hash:     []byte{}}
		blockHash       = block.calcHashForBlock(nonce)
		blockHashAsInt  = binary.LittleEndian.Uint32(blockHash)
		nonce           = nonce + 1
	}

	block.Hash    = blockHash
	blockWrapper := &BlockWrapper{Block: block, Sender: n.address}

	blockWrapperChannel <- blockWrapper
    mineBlock(blockWrapperChannel, n)
}






