package main 

/*
Every communication between nodes over the TCP connection will be 
through sending this communication object.  The ID represents the type:

0 - means we will be receiving a block
1 - means we will be receiving a slice of sent addresses
2 - means we were requested to send conections
3 - means we will be receiving a blockchain
4 - means we were requested to send your blockchain 
*/

type Communication struct {
    ID 				int
    Block        	*BTNode
    SentAddresses   []string
    Blocktree 		*BlockTree
}

// for testing
// func newComm(ID int) Communication{
// 	switch ID {
// 	case 0:
// 		return Communication{0, BlockWrapper{genesisBlock, false, "127.0.0.1:1999"}, []string{}, Blockchain{[]Block{}}}
// 	case 1:
// 		return Communication{1, emptyBlockWrapper(), []string{}, Blockchain{[]Block{}}}
// 	case 2:
// 		return Communication{2, emptyBlockWrapper(), []string{}, Blockchain{[]Block{}}}
// 	case 3:
// 		return Communication{3, emptyBlockWrapper(), []string{}, Blockchain{[]Block{}}}
// 	case 4:
// 		return Communication{4, emptyBlockWrapper(), []string{}, Blockchain{[]Block{}}}
// 	default:
// 		return Communication{0, BlockWrapper{genesisBlock, false, "127.0.0.1:1999"}, []string{}, Blockchain{[]Block{}}}
// 	}
// }