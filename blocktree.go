package main 

import(
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

type BTNode struct {
	Height	   uint32
	Parent     *BTNode
	ParentHash []byte
	Data       string
	Hash 	   []byte
}

type BlockTree struct {
	Levels [][]*BTNode
}

var genesisNode = BTNode{Height: 0, Parent: nil, ParentHash: []byte{0}, Data: "Genesis", Hash: []byte{0}}

func (bt *BlockTree) addBTNodeIfValid(newBTNode *BTNode) {
	parentHeight         := newBTNode.Height - 1
	nodesAtLevelOfParent := bt.Levels[parentHeight]

	for _ , oldBTNode := range nodesAtLevelOfParent {
		if oldBTNode.isValidNextBTNode(newBTNode) {
			// append to parent level
			if uint32(len(bt.Levels)) <= newBTNode.Height { // we define genesis block at height 0
				// this is now the longest chain, append a new level
				var newLevel []*BTNode
				newLevel = append(newLevel, newBTNode) // new level containing the only block high enough
				bt.Levels = append(bt.Levels, newLevel) // should automatically be at correct height
			} else {
				// not the longest chain, directly inject into height at newBTNode.height		
				bt.Levels[newBTNode.Height] = append(bt.Levels[newBTNode.Height], newBTNode)
			}
		} else {
			fmt.Println("No matching node found")
		}
	}
} // should check to see which is now the longest

func (oldBTNode *BTNode) isValidNextBTNode(newBTNode *BTNode) bool {
	heightValid := oldBTNode.Height + 1 == newBTNode.Height
	parentValid := (oldBTNode == newBTNode.Parent)
	hashValid   := testEqByteSlice(oldBTNode.Hash, newBTNode.ParentHash)

	return heightValid && parentValid && hashValid
}

func equalBTNodes(b1, b2 BTNode) bool {
	heightEq     := b1.Height == b2.Height
	parentHashEq := testEqByteSlice(b1.ParentHash, b2.ParentHash)
	dataEq       := b1.Data == b2.Data
	hashEq 		 := testEqByteSlice(b1.Hash, b2.Hash)

	return heightEq && parentHashEq && dataEq && hashEq
}

func (b *BTNode) calcBTNodeHash(){
	height := make([]byte, 4)
	binary.LittleEndian.PutUint32(height, b.Height)
	data := []byte(b.Data)

	h := sha256.New()
	h.Write(height)
	h.Write(data)
	h.Write(b.ParentHash)

	b.Hash = h.Sum(nil)
}
