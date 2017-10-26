package main

/*
Every communication between nodes over the TCP connection will be
through sending this communication object.  The ID represents the type:

0 - means we will be receiving a BlockWrapper
1 - means we will be receiving a slice of sent addresses
2 - means we were requested to send conections
3 - means we will be receiving a blockchain
4 - means we were requested to send your blockchain
*/

type Communication struct {
	ID            int
	BlockWrapper  BlockWrapper
	SentAddresses []string
	Blockchain    Blockchain
	Packet        Packet
}
