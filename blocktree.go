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

type BlockWrapper struct {
	Block  *BTNode
	Sender string
}

type BlockTree struct {
	Levels [][]*BTNode
	Top        *BTNode
}

var genesisNode = BTNode{Height: 0, Parent: nil, ParentHash: []byte{0}, Data: "Genesis", Hash: []byte{0}}

func emptyBlock() BTNode {
	return BTNode{Height: 0, Parent: nil, ParentHash: nil, Data: "", Hash: nil}
}

/* 
addBTNodeIfValid takes a proposed block and an existing block tree
and checks if the parent of this node exists on the blockchain. If
it does, it adds it to the appropirate level of the blocktree and
returns true, otherwise returns false. 
*/
func (bt *BlockTree) addBTNodeIfValid(newBTNode *BTNode) bool {
	parentHeight         := newBTNode.Height - 1
	if uint32(len(bt.Levels)) <= parentHeight  {
		// does not have the parent, so appears as nil.
		return false
	}
	nodesAtLevelOfParent := bt.Levels[parentHeight]

	for _ , oldBTNode := range nodesAtLevelOfParent {
		if oldBTNode.isValidNextBTNode(newBTNode) {
			// append to parent level
			if uint32(len(bt.Levels)) <= newBTNode.Height { // we define genesis block at height 0
				// this is now the longest chain, append a new level
				var newLevel []*BTNode
				newLevel = append(newLevel, newBTNode) // new level containing the only block high enough
				bt.Levels = append(bt.Levels, newLevel) // should automatically be at correct height
				bt.Top = newBTNode
				fmt.Println("just added new top of blockchain")
			} else {
				// not the longest chain, directly inject into height at newBTNode.height		
				bt.Levels[newBTNode.Height] = append(bt.Levels[newBTNode.Height], newBTNode)
			}
			return true
		} else {
			fmt.Println("No matching node found")
		}
	}
	return false
} // should check to see which is now the longest

func (oldBTNode *BTNode) isValidNextBTNode(newBTNode *BTNode) bool {
	heightValid    := oldBTNode.Height + 1 == newBTNode.Height
	// fmt.Printf("Height valid: %v\n", heightValid)
	var parentValid bool
	if newBTNode.Parent != nil{
		parentValid = equalBTNodes(*oldBTNode, *newBTNode.Parent)		
	} else {
		parentValid = false
	}
	// fmt.Printf("Parent valid: %v\n", parentValid)
	parentHashValid := testEqByteSlice(oldBTNode.Hash, newBTNode.ParentHash)
	// fmt.Printf("Parent hash valid: %v\n", parentHashValid)

	/* 
	need to include hash valid that checks if the hash of this
	block (the thing added by calcBTNodeHash()) is correct.

	evenutally I will need to expand the calcBTNodeHash() function to
	include trandactions for the cryptocoin branch of this project
	*/

	return heightValid && parentValid && parentHashValid
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

/*
This is the function that decides which branch in the blocktree is most valid.
Currently it is set to the longest chain, but could abide by other rules.
*/