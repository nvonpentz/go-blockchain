package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	// "github.com/nvonpentz/go-hashable-keys"
)

/*
control.go handles all user interaction with the node
done by entering text via commandline
*/

func listenForUserInput(blockWrapperChannel chan *BlockWrapper, packetChannel chan Packet, n *Node) {
	reader := bufio.NewReader(os.Stdin) //constantly be reading in from std in
	input, err := reader.ReadString('\n')
	if err != nil || input == "\n" {
	} else {
		fmt.Println()
		go handleUserInput(input, blockWrapperChannel, packetChannel, n)
	}
}

func handleUserInput(input string, blockWrapperChannel chan *BlockWrapper, packetChannel chan Packet, n *Node) {
	outgoingArgs := strings.Fields(strings.Split(input, "\n")[0]) // remove newline char and seperate into array by whitespace
	arg0 := strings.ToLower(outgoingArgs[0])
	switch arg0 {
	case "mine":
		go mineBlock(blockWrapperChannel, n)
		listenForUserInput(blockWrapperChannel, packetChannel, n)
	case "getchain":
		if n.seed == "" {
			fmt.Println("You must have a seed node to request a blockchain")
		} else {
			seedConn := n.getConnForAddress(n.seed)
			requestBlockchain(seedConn)
		}
		fmt.Println()
		listenForUserInput(blockWrapperChannel, packetChannel, n)

	case "getconns":
		if n.hasConnectionOfAddress(n.seed) {
			seedConn := n.getConnForAddress(n.seed)
			fmt.Println("Requesting more connections from seed " + n.seed + " ...")
			requestConnections(seedConn)
		} else {
			fmt.Println("You are not connected to your seed node to make a request..")
		}
		fmt.Println()
		listenForUserInput(blockWrapperChannel, packetChannel, n)
	case "node":
		n.printNode()
		fmt.Println()
		listenForUserInput(blockWrapperChannel, packetChannel, n)
	case "genkeys":
		keys := GenerateNewKeypair()
		fmt.Printf("Public: %v\nPrivate: %v\n", string(keys.Public), string(keys.Private))
		fmt.Println()
		listenForUserInput(blockWrapperChannel, packetChannel, n)
	case "upload":
		reader := bufio.NewReader(os.Stdin) //constantly be reading in from std in

		// ask for file name
		fmt.Println("Enter the name of file you wish to save on the blockchain")
		filePath, err := reader.ReadString('\n')
		if err != nil || filePath == "\n" {
			fmt.Println(err)
			fmt.Println("Please enter a valid filepath. Enter 'upload' to begin again.")
			listenForUserInput(blockWrapperChannel, packetChannel, n)
			break
		}
		filePath = strings.Trim(filePath, "\n")

		// ask for public key
		fmt.Println("Enter your public key to associate with the document")
		publicKey, err := reader.ReadString('\n')
		if err != nil || publicKey == "\n" {
			fmt.Println(err)
			listenForUserInput(blockWrapperChannel, packetChannel, n)
		}
		publicKey = strings.Trim(publicKey, "\n")

		// ask for private key
		fmt.Println("Enter your private key to sign the document")
		privateKey, err := reader.ReadString('\n')
		if err != nil || privateKey == "\n" {
			fmt.Println(err)
		}

		// reconstruct keypair
		privateKey = strings.Trim(privateKey, "\n")
		keyPair := Keypair{Public: []byte(publicKey), Private: []byte(privateKey)}

		// create packet and print packet hash to user
		packet := createPacket(filePath, keyPair)
		packetHashHex := hex.EncodeToString(packet.Hash)
		fmt.Printf("This the hash of your packet: %v\n", packetHashHex)

		// check validity of package
		if verifyPacketSignature(packet) {
			fmt.Println("Your packet is valid, sending out to network!")

			// send to packet channel
			packetChannel <- packet
		} else {
			fmt.Println("Your packet was invalid, and will not be sent to blockchain")
		}

		// go back user input as normal
		listenForUserInput(blockWrapperChannel, packetChannel, n)
	case "lookup":
		reader := bufio.NewReader(os.Stdin) //constantly be reading in from std in

		// ask for packet hash
		fmt.Println("Enter the hash (in hex) you wish to lookup on the blockchain")
		packetHash, err := reader.ReadString('\n')
		if err != nil || packetHash == "\n" {
			fmt.Println(err)
			fmt.Println("Please enter a valid filepath. Enter 'lookup' to begin again.")
			listenForUserInput(blockWrapperChannel, packetChannel, n)
			break
		}
		packetHash = strings.Trim(packetHash, "\n")
		packetHashBytes, err := hex.DecodeString(packetHash)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Please enter valid packetHash. Enter 'lookup' to begin again.")
			listenForUserInput(blockWrapperChannel, packetChannel, n)
			break
		}

		// ask for public key
		fmt.Println("Enter the public key associated with that hash")
		publicKey, err := reader.ReadString('\n')
		if err != nil || publicKey == "\n" {
			fmt.Println(err)
			fmt.Println("Please enter a valid publicKey. Enter 'lookup' to begin again.")
			listenForUserInput(blockWrapperChannel, packetChannel, n)
			break
		}
		publicKey = strings.Trim(publicKey, "\n")
		publicKeyBytes := []byte(publicKey)

		// return whether or not the public key validates this packet hash
		packet := n.blockchain.findPacketByHashAndPublicKey(packetHashBytes, publicKeyBytes)
		emptyPacket := Packet{}
		if equalPackets(packet, emptyPacket) {
			fmt.Println("The document you requested does not exist with that combination hash and public key")
		} else {
			fmt.Println("Found the packet you were looking for:")
			fmt.Println(packet)
		}
		listenForUserInput(blockWrapperChannel, packetChannel, n)
	case "help":
		showNodeHelp()
		fmt.Println()
		listenForUserInput(blockWrapperChannel, packetChannel, n)
	default:
		fmt.Println("Enter 'help' for options.")
		listenForUserInput(blockWrapperChannel, packetChannel, n)
	}
}
