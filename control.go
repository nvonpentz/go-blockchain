package main

import (
	"bufio"
	"os"
	"fmt"
	"strings"
)

/*
control.go handles all user interaction with the node

*/

func listenForUserInput(minedBlockChannel chan Block, blockWrapperChannel chan *BlockWrapper, n *Node) {
    for {
        reader := bufio.NewReader(os.Stdin) //constantly be reading in from std in
        input, err := reader.ReadString('\n')
        if (err != nil || input == "\n") {
        } else {
            fmt.Println()
            go handleUserInput(input, blockWrapperChannel, n)
        }
    }
}

// func listenToUserInputChannel(userInputChannel    chan string,
//                               minedBlockChannel   chan Block,
//                               blockWrapperChannel chan *BlockWrapper,
//                               myNode              *Node) {
// 	for {
//         input := <- userInputChannel // user entered some input
//         handleUserInput(input, minedBlockChannel, blockWrapperChannel, myNode)
// 	}
// }

func handleUserInput(input string, blockWrapperChannel chan *BlockWrapper, n *Node) {
    fmt.Println("handling user input")
    outgoingArgs := strings.Fields(strings.Split(input,"\n")[0]) // remove newline char and seperate into array by whitespace
    arg0 := strings.ToLower(outgoingArgs[0])
    switch arg0 {
    case "mine":
        // go listenToMinedBlockChannel(minedBlockChannel, blockWrapperChannel, n)
        go mineBlock(&n.blockchain, blockWrapperChannel, n)                        
    case "getchain":
        if n.seed == "" {
            fmt.Println("You must have a seed node to request a blockchain")
        } else{
            seedConn := n.getConnForAddress(n.seed)
            n.requestBlockchain(seedConn)                        
        }
    case "getconns":
        if n.hasConnectionOfAddress(n.seed){
            seedConn := n.getConnForAddress(n.seed)
            fmt.Println("Requesting more connections from seed " + n.seed + " ...")
            n.requestConnections(seedConn)
        } else {
            fmt.Println("You are not connected to your seed node to make a request..")
        }
    case "node":
        n.printNode()
    case "help":
        showNodeHelp()
    default:
        fmt.Println("Enter 'help' for options.")
    }
}