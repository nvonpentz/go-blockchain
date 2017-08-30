# go-blockchain
Prove you have an idea before anyone else by creating a hash of the document of your thoery, digitally signing it, and uploading it to the blockchain.

Suppose you have an interesting new theory about the world, and you want to be able to prove you had this idea, but don't want to publish it yet because the theory is unfinished.  You can acheive this by typing the theory in a document, creating a hash of this document, and signing the hash with my private key.  You combine the document hash, signature, and my public key into a single `packet` of information, and upload it to the blockchain.

Then, if someone else comes along and claims they had this idea first, you can demonstrate that you was first by making my document public.  Now anyone can hash my document and check if to see if the hash its on the blockchain.  They can also verify that it is indeed my public key which created the signature, and thus you am the only one able to create.

Now timestamping on a blockchain is actually more difficult than you'd hope since there is no notion of global time.  The solution is to include a sample from a current news article (or some document that proves you wrote it at the same time) at the end of the document you wish to upload to the chain.

## Usage
```
NAME:
   go-blockchain

USAGE:
   go-blockchain [global options]

COMMANDS:
   go-blockchain      launches a node

GLOBAL OPTIONS:
    -l, --listen     assigns the listening port for the server        (default = 1999).
    -s, --seed       assigns the port of the seed                     (default = 2000).
    -p, --public     launch node using a your public IP               (default = false).
    -h, --help       prints help information

NODE COMMANDS:
    getconns  requests the list of nodes from your seed node and attempts to connect to each
    getchain  requests seed node for their version of the blockchain
    genkeys   generates and prints a public and private keypair
    node      prints the data associated with your node
    upload    initates the process of uploading a signed document hash to the blockchain
    lookup    initates the process of verifying a document hash and public keypair is on the blockchain
    help      prints the node command help information
```
## Getting started
### Setup
To install, in Terminal, `cd` into your directory containing Go projects and enter:
```
git clone https://github.com/nvonpentz/go-blockchain.git
go build
go install
```
### Starting a node
To launch the first node on the server:
```
go-blockchain -l 1999
```
This launches the node and tells it to listen for incoming connections on port `:1999`.  Let's launch another node and connect them.  In a separate terminal window:
```
go-blockchain -l 2000 -s 1999
```
This launches another node, and specifies the seed node to be at port `:1999` and listen port to be `:2000`.  The nodes will connect. (Note:  The default listening port is `:1999` but in order to simulate the network on a single computer, we listen on different ports.)

### Requesting new connections
If you want to connect to more than just your seed node, you can request a list of addresses from your seed node, and atttempt to connect to each with `getconns`. To test this out, start three nodes on the network such that each node is connected to only one peer:

```
Node1 connected to Node2
Node2 connected to Node3
```

From Node3, enter `getconns`.  It will ask Node2 for its connections, and Node2 will send back a list of addresses to connect to.  Node3 will then connect to those that it isn't connected with already (ie Node1).  Enter `node` to see the connections, and you will see that it is connected to both Node2 and Node1.

### Upload a document

### Start mining
After you have booted up the node, enter `mine`, and your node will 
attempt to solve the mining puzzle to mine a block.  Once a valid nonce is found your node will automatically send them to the network when they are mined.

## Code Explanation
Every node is considered a **full node** and can:
* Mine new blocks, add them to their blockchain, and send to connected nodes
* Receive blocks mined by other nodes, check if they are valid, and send them all their connections if it is
* Create packets and forward to the network
* Accept valid packets and add them to their current block for which they are trying to solve the mining puzzle
* Verify that a valid packet is exists already on the network

### Mining
Blocks are mined finding a nonce value such that:

  SHA256(block index ᛫ previous block hash ᛫ block packets hash ᛫ nonce) < difficulty target

* The difficulty target, 200, is hard coded in `node.go`.
* The mining algorithm is found in `mine.go`
* The block hashing function is found in `block.go`.

Once you mine a block, we create a new struct, a `blockWrapper` and send it to your node's the blockWrapper channel where all blocks (including block sent from the network) are processed.  A block wrapper consists of the original block, as well the most recent sender.

```go
type BlockWrapper struct {
    Block  Block
    Sender string
}
```

### Validating blocks
When a block is sent to your node's blockchannel (either from successfully finding a nonce, or from receiving a block from one of your peers) your node checks to see if it has seen the block before.  If it hasn't it checks the validity of the block.  A block is valid if:

* It's index is one greater than the previous block
* Its previous hash is equal to the previous block's hash
* All the signatures in the block's list of packets are valid
* The hash of the block computed by your computer matches the claimed hash on the block
* The hash of the block is below difficulty target

If the block is valid, it is added to the of seen blocks, and forwards it to all of its connections.  Blocks are validated in the `isValidNextBlock` function in `block.go`.

There is a special circumstance in which a valid block is sent to your node, but your node does not recognize it as valid, because this blocks index is more than one ahead than the block at the tip of your node's blockchain.  This creates a bad scenario in which your node will mark the block as invalid, and add it to it's list of seen blocks.  So even if you were to eventually receive intermediate blocks between your node's tip and this block, your node would could never assimilate it, as it has discarded the block.

The solution used in this blockchain is to send a request for the entire blockchain to the node who sent a block whose index is more than one greater than your nodes highest block.  In this case, your node will validate the entire chain, and if it is all valid, replace its current chain with the one received from its peers.  This is why the `Sender` field is included in the `blockWrapper`, in order to request entire blockchains from nodes who send a block which appears to be invalid, but might be valid in the context of the sending node's blockchain.

### Validating packets
Packets are the atomic element of this blockchain.

### Network
Nodes communicate via TCP.  Every communication passed between nodes in the network is actually just a instance of `Communication` struct:
```
type Communication struct {
    ID            int
    BlockWrapper  BlockWrapper // wrapper for sent block
    SentAddresses []string
    Blockchain    Blockchain
    Packet        Packet
}
```
Depending on the value of the `Communication.ID`, the communication instance is either:
* A block (ID = 0)
* A response of connection addresses (ID = 1)
* A request for connection addresses (ID = 2)
* A response of a blockchain (ID = 3)
* A request to send a blockchain (ID = 4)
* A packet (ID = 5)

When a communication is sent over the network, it is parsed by the `listenToConnection()` go routine, and redirects the Datarmation to the appropriate channel.

## Improvements
* change myNode listen to connections to not be a function of my node or atleast get rid of the n.connections argument CHECK
* make the hashes base 58 so they can be human readable when strings. CHECK
* be able to check for a specific transaction hash CHECK
* add channels to miner so it can be updated with new blocks/new packets as they come


