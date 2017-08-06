package main

import(
    "fmt"
)

func showGlobalHelp() {
fmt.Println(
`NAME:
   go-blockchain - blockchain network

USAGE:
   go-blockchain [global options]

COMMANDS:
   go-blockchain      launches a node

GLOBAL OPTIONS:
    -l, --listen     assigns the listening port for the server        (default = 1999).
    -s, --seed       assigns the port of the seed                     (default = 2000).
    -j, --join       attempt to join the network                      (default = false).
    -h, --help       prints this help info

NODE COMMANDS:
    send      sends the subsequent text to the network
    getconns   requests the list of nodes from your seed node and attempts to connect to each
    getchain  requests seed node for their version of the blockchain
    node      prints the information associated with your node
    help      prints the node command help info`)
}

func showNodeHelp(){
fmt.Println(
`
NODE COMMANDS:
    send      sends the subsequent text to the network
    getconns  requests the list of nodes from your seed node and attempts to connect to each
    node      prints the information associated with your node
    help      prints the node command help info`)
}