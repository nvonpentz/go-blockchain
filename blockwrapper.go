package main 

/*
A BlockWrapper is simply a message that will be sent throughout the network.  
It includes the actual message as well as the addresses of nodes who have
already received the BlockWrapper
*/

// type BlockWrapper struct {
//     Block Block
//     BeenSent bool
//     Sender string
// }

// func (t *BlockWrapper) updateBeenSent() {
//     t.BeenSent = true
// }

// func (t *BlockWrapper) updateSender(address string){
//     t.Sender = address
// }

// // testing
// // func emptyBlockWrapper() BlockWrapper{
// // 	return BlockWrapper{emptyBlock(), false, "127.0.0.1"}
// // }