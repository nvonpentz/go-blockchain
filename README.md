# go-p2p

A simple implementation of a peer-to-peer network over TCP, written to be the basis of a simple blockchain network.

Once connected to peer nodes, a node can send information (a string text) to each node on the network.  Each node will save that information for itself and pass the message along to it's peers who haven't seen the information yet.

## Usage
```
NAME:
   go-p2p - peer to peer network

USAGE:
   go-p2p [global options]

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
git clone https://github.com/nvonpentz/go-p2p.git
go build
go install
```
### Starting A Node
To launch the first node on the server:
```
go-p2p -l 1999
```
This launches the node and tells it to listen for incoming connections on port `:1999`.  Let's launch another node and connect them.  In a separate terminal window:
```
go-p2p -l 2000 -s 1999
```
This launches another node, and specifies the seed node to be at port `:1999` and listen port to be `:2000`.  The nodes will connect. (Note:  The default listening port is `:1999` but in order to simulate the network on a single computer, we listen on different ports.)

### Send Information To Network
To pass information between nodes simply type `send` followed by the information you wish to log.  This will send the information across the entire newtork, and each node will append it in a slice, `myNode.information`.

### Requesting New Connections
If you want to connect to more than just your seed node, you can request a list of addresses from your seed node, and atttempt to connect to each with `request`. To test this out, start three nodes on the network such that each node is connected to only one peer:

```
Node1 connected to Node2
Node2 connected to Node3
```

From Node3, enter `request`.  It will ask Node2 for its connections, and connect to those that it isn't connected with already (ie Node1).  Enter `node` to see the connections, and you will see that it is connected to both Node2 and Node1.

### Help
To review a list of commands available to the node enter `help` once you have launched a node.  To see the commands available to the program as a whole enter `go-p2p -h` before you have launched the node.








