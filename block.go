package main

import ("fmt"
		"time"
		// "math/rand"
		"crypto/sha256"
		"encoding/binary"
		)

type Blockchain struct {
	Blocks []Block
}

type Block struct {
	Index uint32
	PrevHash []byte
	Info string
	Hash []byte
}

var genesisBlock = Block{0, []byte{0}, "genesis transaction", []byte{0}}

func (blockchain *Blockchain) addBlock(block Block) {
	if blockchain.isValidBlock(block) == true {
		blockchain.Blocks = append(blockchain.Blocks, block)
		// fmt.Printf("Added block %d to blockchain \n", block.Index)
	} else {
		fmt.Println("Did not add block to blockchain")
	}
}

func (blockchain *Blockchain) isValidBlock(block Block) bool{
	lastBlockInChain := blockchain.getLastBlock()
	return isValidBlocks(lastBlockInChain, block)
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

func (blockchain *Blockchain) mineBlock(blockChannel chan Block){
	fmt.Println("-> begin mine")
	// sleepTime := time.Duration((rand.Int() % 10) + 5) //use randomness for now
    time.Sleep(time.Second * 2)

    prevBlock     := blockchain.getLastBlock()
	newBlockIndex := prevBlock.Index + 1
	newBlockInfo  := "new block!"
	newBlock := Block{newBlockIndex, prevBlock.Hash, newBlockInfo, []byte{}}

	newBlockHash := calcHashForBlock(newBlock)
	newBlock      = Block{newBlockIndex, prevBlock.Hash, newBlockInfo, newBlockHash}

	blockChannel <- newBlock
}

func calcHashForBlock(block Block) []byte {
	blockHash := sha256.New()

	nbIndexBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(nbIndexBytes, block.Index)
	pbHashBytes  := []byte(block.PrevHash)
	nbInfoBytes  := []byte(block.Info)
	toHash := append(nbIndexBytes, pbHashBytes...)
	toHash  = append(toHash, nbInfoBytes...)
	blockHash.Write(toHash)

	return blockHash.Sum(nil)
}

func isValidBlocks(b1 Block, b2 Block) (bool){
	isValidIndex := b2.Index == b1.Index + 1
	/*
	incoming blocks previous hash has to equal the hash of the 
	*/
	return isValidIndex
}


