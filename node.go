package main

import(
    "fmt"
    "net"
    "os"
    "encoding/gob"
    "strings"
    "regexp"
    "net/http"
    "io/ioutil"
)

type Node struct {
    connections map[net.Conn]int
    nextConnID  int
    blocktree   *BlockTree
    address     string
    seed        string
    seenBlocks  map[string]bool
}

/*-----------*
 *  METHODS  * 
 *-----------*/

func (myNode Node) run(listenPort string, seedInfo string, publicFlag bool) {
    joinFlag := false
    if seedInfo != "" { joinFlag = true }
    
    // specify ports to seed and listen to
    myNode.updatePorts(listenPort, seedInfo, publicFlag)

    // create channels
    // blockWrapperChannel      := make(chan *BlockWrapper)
    blockChannel             := make(chan *BlockWrapper)
    newConnChannel           := make(chan net.Conn) // new connections added
    disconChannel            := make(chan net.Conn) // new disconnestion
    connRequestChannel       := make(chan net.Conn) // received a request to send connections 
    sentAddressesChannel     := make(chan []string) // received addresses to make connections
    // minedBlockChannel        := make(chan Block)    // new block was mined
    blockchainRequestChannel := make(chan net.Conn)
    sentBlockTreeChannel     := make(chan *BlockTree)

    // listen to user input
    go listenForUserInput(blockChannel, &myNode)

    // listen on network
    listenForConnections(listenPort, newConnChannel)
    if joinFlag { // if the user requested to join a seed node // need to make sure you can't join if you don't supply a seed
        fmt.Println("Dialing seed node at port " + seedInfo + "...")
         go dialNode(myNode.seed, newConnChannel)
    }

    myNode.printNode()

    // handle network go routines
    for {
        select {
            case conn         := <- newConnChannel: // listener picked up new conn
                myNode.nextConnID = myNode.nextConnID + 1
                myNode.connections[conn] = myNode.nextConnID // assign connection an ID
                go listenToConn(conn, blockChannel, disconChannel, connRequestChannel, sentAddressesChannel, blockchainRequestChannel, sentBlockTreeChannel)

            case discon       := <- disconChannel: // established connection disconnected
                connID := myNode.connections[discon]
                delete(myNode.connections, discon) // remove the connection from the nodes list of connections
                fmt.Printf("* Connection %v has been disconnected \n", connID)
            case blockWrapper := <- blockChannel:
                block  := blockWrapper.Block
                seenBlock := myNode.seenBlocks[string(block.Hash)] == true
                if !seenBlock {
                    myNode.seenBlocks[string(block.Hash)] = true
                    myNode.blocktree.addBTNodeIfValid(block)
                    blockValid := myNode.blocktree.addBTNodeIfValid(block)
                    if blockValid {
                        myNode.forwardBlockToNetwork(&BlockWrapper{Block: block, Sender: myNode.address}, myNode.connections)                        
                        fmt.Println("sent blockchain to network")
                    } else {
                        fmt.Println("block was not considered valid, making request for whole chain to compare..")                        
                        requestBlockchain(myNode.getConnForAddress(blockWrapper.Sender), block) //request blockchain ending in block, ba
                    }
                }
            // case blockWrapper := <- blockWrapperChannel:  // new blockWrapper sent to node // handles adding, validating, and sending blocks to network
            //     myNode.handleBlockWrapper(blockWrapper)
            case conn         := <-  connRequestChannel:  // was requested addresses to send
                addressesToSendTo := myNode.getRemoteAddresses()
                sendConnectionsToNode(conn, addressesToSendTo)

            case addresses    := <- sentAddressesChannel:  //received addresses to add
                fmt.Printf("Seed node sent these addresses to connect to:\n-> %v\n", addresses)
                myNode.handleSentAddresses(addresses, newConnChannel)

            // case conn         := <- blockchainRequestChannel:
            //     sendBlockchainToNode(conn, myNode.blockchain)

            // case blockchain   := <- sentBlockTreeChannel: // node was sent a blockchain
            //     myNode.handleSentBlockchain(blockchain)
        }

    }
}

func (n *Node) updatePorts(listenPort string, seedInfo string, publicFlag bool) {
    if publicFlag{
        n.seed = seedInfo + ":1999" // if public ip, seed is specifiec seedInfo:1999
        n.address = getPublicIP() + ":1999" // must set up port forwarding
    } else { 
        n.seed = getPrivateIP() + ":" + seedInfo  // no default seed
        n.address = getPrivateIP() + listenPort
    }
}

func (n *Node) forwardBlockToNetwork(blockWrapper *BlockWrapper, connections map[net.Conn]int) {
    for conn, _ := range connections { // loop through all this nodes connections
        // destinationAddr := conn.RemoteAddr().String() // get the destination of the connection
        communication := Communication{0, blockWrapper, []string{}, []*BTNode{}}
        encoder       := gob.NewEncoder(conn)
        encoder.Encode(communication)        
    }
}


// func (n *Node) forwardBlockWrapperToNetwork(blockWrapper BlockWrapper, connections map[net.Conn]int) {
//     for conn, _ := range connections { // loop through all this nodes connections
//         // destinationAddr := conn.RemoteAddr().String() // get the destination of the connection
//         communication := Communication{0, blockWrapper, []string{}}
//         encoder       := gob.NewEncoder(conn)
//         encoder.Encode(communication)        
//     }
// }

// func (n *Node) handleBlockWrapper(blockWrapper *BlockWrapper){
    
//     notMinedAndValid   := n.seenBlocks[string(blockWrapper.Block.Hash)] == false  && blockWrapper.BeenSent == true && n.blockchain.isValidBlock(blockWrapper.Block)
//     notMinedAndInvalid := n.seenBlocks[string(blockWrapper.Block.Hash)] == false  && blockWrapper.BeenSent == true && !n.blockchain.isValidBlock(blockWrapper.Block)
//     minedButNotSent    := n.seenBlocks[string(blockWrapper.Block.Hash)] == true   && blockWrapper.BeenSent == false
//     alreadySent        := n.seenBlocks[string(blockWrapper.Block.Hash)] == true   && blockWrapper.BeenSent == true
    
//     if notMinedAndValid {
//         n.seenBlocks[string(blockWrapper.Block.Hash)] = true
//         n.blockchain.addBlock(blockWrapper.Block)
//         fmt.Printf("[notMinedAndValid] Added block #%v sent from network to my blockchain, and sending it to network\n", blockWrapper.Block.Index)
//         blockWrapper.updateSender(n.address)
//         n.forwardBlockWrapperToNetwork(*blockWrapper, n.connections) // forward messages to the rest of network
//     } else if notMinedAndInvalid {
//         n.seenBlocks[string(blockWrapper.Block.Hash)] = true
//         myBlockchainLength := n.blockchain.getLastBlock().Index
//         if blockWrapper.Block.Index > myBlockchainLength {
//             connThatSentHigherBlockIndex := n.getConnForAddress(blockWrapper.Sender)
//             fmt.Println("I was sent a block with a higher index, now requesting full chain to validate")
//             requestBlockchain(connThatSentHigherBlockIndex)
//         }
//         fmt.Printf("[notMinedAndInvalid] Did not add block #%v sent from network to my chain, did not forward\n", blockWrapper.Block.Index)
//     } else if minedButNotSent { //mined but not sent out yet,
//         blockWrapper.updateBeenSent()
//         blockWrapper.updateSender(n.address) 
//         fmt.Printf("[minedButNotSent] Mined block #%v, sending to network\n", blockWrapper.Block.Index)
//         n.forwardBlockWrapperToNetwork(*blockWrapper, n.connections) // forward messages to the rest of network
//     } else if alreadySent{
//         fmt.Printf("[alreadySent] Already seen block #%v, did not forward\n", blockWrapper.Block.Index)
//     } else {
//         fmt.Println("Some other case, this should not occur:")
//     }
// }

func (n *Node) handleSentAddresses(addresses []string, newConnChannel chan net.Conn){
    approvedAddresses := []string{}
    for i := range addresses {
        r, _ := regexp.Compile(":.*") // match everything after the colon
        port := r.FindString(addresses[i])
        if len(port) == 5 {  // in a real network this should simply be 1999
            go dialNode(addresses[i], newConnChannel)
            approvedAddresses = append(approvedAddresses, addresses[i])
        }
    }
    fmt.Printf("These new connections will be added:\n->%v\n", approvedAddresses)
}

// func (n *Node) handleSentBlockchain(blockchain Blockchain){
//     fmt.Println("You were sent a blockchain")
//     if blockchain.isValidChain() {
//         n.blockchain = blockchain
//         fmt.Println("Blockchain accepted: ")
//         fmt.Println(blockchain)
//     } else {
//         fmt.Println("Blockchain rejected, invalid")
//     }
// }

func (n Node) getConnForAddress(address string) (net.Conn){
    for conn := range n.connections {
        remoteAddr := conn.RemoteAddr().String()
        if address == remoteAddr {
            return conn
        }
    }
    var emptyConn net.Conn
    return emptyConn
}

func (n Node) getRemoteAddresses() (remoteAddresses []string) {
    for conn, _ := range n.connections {
        remoteAddr := conn.RemoteAddr().String()
        remoteAddresses = append(remoteAddresses, remoteAddr)
    }
    return remoteAddresses
}

// this could be made more efficient by not using getRemoteAddresses()
func (n Node) hasConnectionOfAddress(address string) (bool) {
    remoteAddresses := n.getRemoteAddresses()
    for i := 0; i < len(remoteAddresses); i++ {
        if address == remoteAddresses[i] {
            return true
        }
    }
    return false
}

// func (n Node) printSeenBlockWrapper(){
//     for blockHashString, _  := range n.seenBlocks{
//         blockHashBytes := []byte(blockHashString)
//         fmt.Printf("  %v\n", blockHashBytes)
//     }
// }

// func (n Node) printBlockchain(){
//     for i := range n.blockchain.Blocks {
//         block := n.blockchain.Blocks[i]
//         fmt.Printf("  Block %d is: \n   PrevHash: %v \n   Info:     %v \n   Hash:     %v \n", i, block.PrevHash, block.Info, block.Hash)
//     }
// }


func (n Node) printConnections(){
    for conn, id := range n.connections {
        localAddr := conn.LocalAddr().String()
        remoteAddr := conn.RemoteAddr().String()
        fmt.Printf(" ID: %v, Connection: %v to %v \n", id, localAddr, remoteAddr)
    }
}

func (n Node) printNode(){
    fmt.Println("*------------------*\nYour Node:\n Connections:")
    fmt.Printf(" Your Address:\n  %v \n Seed Address:\n  %v\n", n.address, n.seed)
    n.printConnections()
    fmt.Println(" Seen Blocks:")
    // n.printSeenBlockWrapper()
    // fmt.Println(" Blockchain:")
    // n.printBlockchain()
    fmt.Println("*------------------*")
}

/*-------------*
 *  FUNCTIONS  * 
 *-------------*/

func newNode() Node {
    myNode := Node{connections: make(map[net.Conn]int),
                   nextConnID: 0,
                   blocktree: &BlockTree{[][]*BTNode{{&genesisNode}}, &genesisNode},
                   address: "",
                   seed: "",
                   seenBlocks: map[string]bool{}}
    return myNode
}

func listenForConnections(port string, newConnChannel chan net.Conn) {
    listener, err := net.Listen("tcp", port)
    if err != nil {
        fmt.Println("There was an error setting up the listener:")
        fmt.Println(err)
    }

    go acceptConn(listener, newConnChannel)
}

func acceptConn(listener net.Listener, newConnChannel chan net.Conn) {
    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Printf("There was an error accepting a new connection: \n %v", err)
        }
        fmt.Println("* Listener accepted connection through port " + conn.LocalAddr().String() + " from " + conn.RemoteAddr().String())
        newConnChannel <- conn //send to conection channel
    }
}

func dialNode(address string, newConnChannel chan net.Conn) {
    conn, err := net.Dial("tcp", address)
    if err != nil {
        fmt.Println("**Make sure there is someone listening at " + address + "**")
        fmt.Println(err)
    }
    fmt.Println("Connection established out of port " + conn.LocalAddr().String() + " dialing to " + conn.RemoteAddr().String())
    newConnChannel <- conn
}

func listenToConn(conn                          net.Conn, 
                  blockChannel             chan *BlockWrapper,
                  disconChannel            chan net.Conn,
                  connRequestChannel       chan net.Conn,
                  sentAddressesChannel     chan []string,
                  blockchainRequestChannel chan net.Conn,
                  sentBlockTreeChannel     chan *BlockTree) {
    for {
        decoder := gob.NewDecoder(conn)
        var communication Communication
        err := decoder.Decode(&communication)

        if err != nil {
            fmt.Println(err)
            break
        }
        switch communication.ID {
        case 0:
            blockChannel <- communication.BlockWrapper
        case 1:
            sentAddressesChannel <- communication.SentAddresses
        case 2:
            fmt.Println("You have been requested to send your connection addresses to a peer at " + conn.RemoteAddr().String() + " ...")
            connRequestChannel <- conn
        // case 3:
        //     sentBlockTreeChannel <- communication.Blocktree
        case 4:
            fmt.Println("You have been requested to send your blockchain address to a peer at " + conn.RemoteAddr().String() + " ...")
            // topOfRequestedBlockchain := communication.BlockWrapper.Block
            // requestedChain := myNode.blocktree.deriveChainFromTopAndTree(topOfRequestedBlockchain)
            // sendRequestedBlockchain(conn, requestedChain)
            // blockchainRequestChannel <- conn
        default:
            fmt.Println("There was a problem decoding the message")
            break
        }
    }
    disconChannel <- conn // disconnect must have occurred if we exit the for loop
}

func sendRequestedBlockchain(conn net.Conn, requestedChain []*BTNode){
    communication := Communication{3, &BlockWrapper{&BTNode{}, ""}, []string{}, requestedChain}
    encoder   := gob.NewEncoder(conn)
    encoder.Encode(communication)
    fmt.Printf("Sent my copy of blockchain to %v", conn.RemoteAddr().String())
}

func requestConnections(conn net.Conn){
    communication := Communication{2, &BlockWrapper{&BTNode{}, ""}, []string{}, []*BTNode{}}
    encoder       := gob.NewEncoder(conn)
    encoder.Encode(communication)
}

func sendConnectionsToNode(conn net.Conn, addresses []string){
    communication := Communication{1, &BlockWrapper{&BTNode{}, ""}, addresses, []*BTNode{}}
    encoder       := gob.NewEncoder(conn)
    encoder.Encode(communication)
}

// func sendBlockchainToNode(conn net.Conn, blockchain Blockchain){
//     communication := Communication{3, &BlockWrapper{&BTNode{}, ""}, []string{}, blockchain}
//     encoder   := gob.NewEncoder(conn)
//     encoder.Encode(communication)
//     fmt.Printf("Sent my copy of blockchain to %v", conn.RemoteAddr().String())
// }

func requestBlockchain(conn net.Conn, proposedBlock *BTNode){
    communication := Communication{4, &BlockWrapper{proposedBlock, ""}, []string{}, []*BTNode{}}
    encoder   := gob.NewEncoder(conn)
    encoder.Encode(communication)
}

// func sendBlockWrapperFromMinedBlock(block Block, blockWrapperChannel chan *BlockWrapper){
//     blockWrapper := BlockWrapper{block, false, ""}
//     blockWrapperChannel <- &blockWrapper
// }

func getPrivateIP() string {
    name, err := os.Hostname()
    if err != nil {
        return ""
    }
    address, err := net.LookupHost(name)
    if err != nil {
        return ""
    }

    return address[0]
}

func getPublicIP() string {
    resp, err := http.Get("http://myexternalip.com/raw")
    if err != nil {
        os.Stderr.WriteString(err.Error())
        os.Stderr.WriteString("\n")
        os.Exit(1)
    }
    defer resp.Body.Close()

    myPublicIP, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println(err)
    }

    myPublicIPstring := strings.Trim(string(myPublicIP), "\n")
    return myPublicIPstring
}

