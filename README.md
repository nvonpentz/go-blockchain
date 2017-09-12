# go-blockchain
Prove you had an idea at a certain date without revealing the idea, by hashing the document explaining your idea, digitally signing it, and uploading it to the blockchain.
---

Suppose you have an interesting new theory and you want to be able to prove you had this idea, but don't want to publish it yet because the theory isn't finished.  You can acheive this by producing a document which explains your idea, creating a hash of it, and signing the hash with your private key.  Combine the document hash, signature, and your public key into a single `packet` of information, and upload it to the blockchain.

Then, if someone claims they had the idea first, you can demonstrate you were first by publishing the document and your public key.  Now anyone can hash your document and verify that it is indeed on the blockchain, and signed by your private key.

Because there is no notion of global time, timestamping on a blockchain is tricky.  The solution is to include a sample from a current news article (or some document that proves you wrote it at the same time) at the end of the document you wish to upload to the chain.

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
go get
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
To upload a document to the blockchain, simply move the file into the go-blockchain directory.  You will need a public/private key pair, so if you don't have any already, enter `genkeys` and new keys will be printed to the terminal.

Once you have your keys, initiate the upload process by entering `upload`.  You will be prompted for the filename, public and private keys.  If the public and private keys match, and the file exists, a `packet` will be created sent forwarded to all your connections

### Verify a document
To verify a document exists on the blockchain, first the hash of the document.

Next, boot up the node and get the most recent version of the blockchain from your seed node, by entering `getchain`.  Now that you have the most up to date blockchain, enter `lookup` to initate the verification process.  Supply the document hash and the public key that claims to own the document, and your node will search through its copy of the blockchain for a valid packet that matches the hash and public key, and return true if a valid packet is found.

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

### Validating packets
Each block is filled with packets. Packets are the atomic elements of this blockchain.

A packet represents a hashed document, a public key (the owner of the document) and a signature over the hash that is verified by the supplied public key:

```go
type Packet struct {
  Hash      []byte
  Signature []byte
  Owner     []byte
}
```

Validating packets is simple; check whether the signature is valid over the hash with the associated public key.

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
* A block                            (ID = 0)
* A response of connection addresses (ID = 1)
* A request for connection addresses (ID = 2)
* A response of a blockchain         (ID = 3)
* A request to send a blockchain     (ID = 4)
* A packet                           (ID = 5)

When a communication is sent over the network, it is parsed by the `listenToConnection()` go routine, and redirects the Datarmation to the appropriate channel.

## Improvements
* add channels to miner so it can be updated with new blocks/new packets as they come instead of recreating a new block for each hash attempt
