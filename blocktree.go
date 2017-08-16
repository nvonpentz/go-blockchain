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
		if oldBTNode.isValidNextBTNode(*newBTNode) {
			// append to parent level
			bt.Levels[newBTNode.Height] = append(bt.Levels[newBTNode.Height], newBTNode)
		} else {
			fmt.Println("No matching node found")
		}
	}
} // should check to see which is now the longest

func (oldBTNode BTNode) isValidNextBTNode(newBTNode BTNode) bool {
	heightValid := oldBTNode.Height + 1 == newBTNode.Height
	fmt.Printf("Height valid: %v\n", heightValid)

	parentValid := (newBTNode.Parent == &oldBTNode)
	fmt.Printf("newBTNode.Parent %v\n", newBTNode.Parent)
	fmt.Printf("&oldBTNode       %v\n", &oldBTNode)
	fmt.Printf("Parent valid:    %v\n", parentValid)

	hashValid   := testEqByteSlice(oldBTNode.Hash, newBTNode.ParentHash)
	fmt.Printf("Hash valid: %v\n", hashValid)

	return heightValid && parentValid && hashValid
}

func equalBTNodes(b1, b2 BTNode) bool {
	heightEq     := b1.Height == b2.Height
	parentHashEq := testEqByteSlice(b1.ParentHash, b2.ParentHash)
	dataEq       := b1.Data == b2.Data
	hashEq 		 := testEqByteSlice(b1.Hash, b2.Hash)

	return heightEq && parentHashEq && dataEq && hashEq
}

func (b *BTNode) calcBTNodeHash() []byte {
	height := make([]byte, 4)
	binary.LittleEndian.PutUint32(height, b.Height)
	data := []byte(b.Data)

	h := sha256.New()
	h.Write(height)
	h.Write(data)
	h.Write(b.ParentHash)

	return h.Sum(nil)
}
