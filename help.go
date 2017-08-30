package main

import(
    "fmt"
)

func showGlobalHelp() {
fmt.Println(
`NAME:
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
    help      prints the node command help information`)
}

func showNodeHelp(){
fmt.Println(
`
NODE COMMANDS:
    getconns  requests the list of nodes from your seed node and attempts to connect to each
    getchain  requests seed node for their version of the blockchain
    genkeys   generates and prints a public and private keypair
    node      prints the data associated with your node
    upload    initates the process of uploading a signed document hash to the blockchain
    lookup    initates the process of verifying a document hash and public keypair is on the blockchain
    help      prints the node command help information`)
}