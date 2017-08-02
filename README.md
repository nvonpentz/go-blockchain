# go-blockchain

A simple implementation of a privatepeer-to-peer blockchain.  There is no cryptocoin associated with this blockchain, it just the data structure.

## Usage
```
NAME:
   go-blockchain - peer to peer network

USAGE:
   go-blockchain [global options]

COMMANDS:
   go-p2p      launches a node

GLOBAL OPTIONS:
    -l, --listen     assigns the listening port for the server        (default = 1999).
    -s, --seed       assigns the port of the seed                     (default = 2000).
    -j, --join       attempt to join the network                      (default = false).
    -h, --help       prints this help info

NODE COMMANDS:
    send      sends the subsequent text to the network
    request   requests the list of nodes from your seed node and attempts to connect to each
    node      prints the information associated with your node
    help      prints the node command help info
```
## Getting Started
### Setup
To install, in Terminal, `cd` into the directory containing Go projects and enter:
```
git clone https://github.com/nvonpentz/go-blockchain.git
go build
go install
```
### Starting A Node
To launch the first node on the server:
```
go-blockchain -l 1999
```
This launches the node and tells it to listen for incoming connections on port `:1999`.  Let's launch another node and connect them.  In a separate terminal window:
```
go-blockchain -l 2000 -s 1999
```
This launches another node, and specifies the seed node to be at port `:1999` and listen port to be `:2000`.  The nodes will connect. (Note:  The default listening port is `:1999` but in order to simulate the network on a single computer, we listen on different ports.)

### Requesting New Connections
If you want to connect to more than just your seed node, you can request a list of addresses from your seed node, and atttempt to connect to each with `getconns`. To test this out, start three nodes on the network such that each node is connected to only one peer:

```
Node1 connected to Node2
Node2 connected to Node3
```

From Node3, enter `getconns`.  It will ask Node2 for its connections, and connect to those that it isn't connected with already (ie Node1).  Enter `node` to see the connections, and you will see that it is connected to both Node2 and Node1.

### Start Mining
To start mining, after you have booted up the node, simpy enter `mine`, and your node will begin mining blocks, and send them to the network when they are mined.  Nodes will automatically validate blocks/blockchains which are sent to them.

## Summary
A node on the network will:
* Mine new blocks, add them to their blockchain, and send to connected nodes
* Receive blocks mined by other nodes, validate them, and send them all their connections

Currently, block are mined every 5 to 10 seconds, and send to the `blockChannel` to be processed.
```{Go}
func (blockchain *Blockchain) mineBlock(blockChannel chan Block){
  fmt.Println("-> begin mining..")

  // sleep between 5 - 10 seconds before mining block to simulate a blockchain
  sleepTime := time.Duration((rand.Int() % 10) + 5)
  time.Sleep(time.Second * sleepTime)

  //create new block
  prevBlock     := blockchain.getLastBlock()
  newBlockIndex := prevBlock.Index + 1
  newBlockInfo  := "new block!"
  newBlock := Block{newBlockIndex, prevBlock.Hash, newBlockInfo, []byte{}}

  // must calculate the hash of this block
  newBlockHash := calcHashForBlock(newBlock)
  newBlock      = Block{newBlockIndex, prevBlock.Hash, newBlockInfo, newBlockHash}

  // send to control center to 
  blockChannel <- newBlock 
}
```
Once in the `blockChannel`, the block is validated by the `isValidBlock()` blockchain method.  A pending block is valid if it's property `PrevHash` is equal to the `Hash` of the current latest block, and if it's `Index` is 1 greater than the latest block's `Index`:
```{Go}
func areValidBlocks(oldBlock Block, newBlock Block) (bool){
  // new block's index must be one greater
  isValidIndex := newBlock.Index == oldBlock.Index + 1

  // new block's previous hash has to equal the hash of the old block
  isValidHash := testEqByteSlice(newBlock.PrevHash, oldBlock.Hash)
  isValidBlock := isValidIndex && isValidHash

  return isValidBlock
}
```
If the block validates with respect to the Node's blockchain, it is sent throughout the network in the form of a `Transmission`.  A transmission is simply the Block, a bool representing whether it has been sent to the network yet, as well as the most recent sender:
```{Go}
type Transmission struct {
    Block Block
    BeenSent bool
    Sender string
}
```
Through the configuration of a transmission, a Node can determine whether it has already seen the transmission and whether a transmission is valid or not

* If the node has already seen the transmission or if the block is not valid, the node does not add the block to it's chain, and does not forward the block to the rest of the network.
* If the node has not seen the transmission, and the block contained within is valid with respect to the nodes current chain, it adds the block to its chain, and forwards it to the rest of the network.
* If a Node receives a transmission containing a block that has a higher `Index` value compared to the last block in its chain (ie. it sees that the node that sent it is claiming to have a **longer chain**), it sends a request to the `Sender` of the transmission for the entire blockchain that is supposedly longer.  Once the Node receives this supposedly longer blockchain, it validates the entire chain, and if it turns out that the chain is valid, the node replaces its own shorter chain with this new valid chain.

This behavior happens within the `transmissionChannel` in `main.go`

## Why Private?
This is a private blockchain, which means it cannot easily be run beyond a private network because of the challenges of getting past routers and NAT.  Theoretically this blockchain would work as public blockchain if users setup portforwarding on their router, or if universal plug and play (UPNP) was implemented.

