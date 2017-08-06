package main

import(
    "fmt"
    "net"
    "bufio"
    "os"
    "encoding/gob"
    "strings"
    "regexp"
)

/*-------------------*
 * STRUCTS & METHODS *
 *-------------------*/

type Node struct {
    connections map[net.Conn]int
    nextConnID int
    blockchain Blockchain
    address string
    seed string
    seenBlocks map[string]bool
}

func (n *Node) incrementConnID() {
    n.nextConnID = n.nextConnID + 1
}

func (n *Node) appendBlock(block Block) {
    n.blockchain.Blocks = append(n.blockchain.Blocks, block)
}

func (n *Node) updateAddress(listenPort string) {
    n.address = getPrivateIP() + listenPort
}

func (n *Node) updateSeed(seedPort string) {
    if seedPort == ":"{ // ie empty seed port **should refactor**
        n.seed = ""
    } else {
        n.seed = getPrivateIP() + seedPort
    }
}

func (n Node) printNode(){
    fmt.Println("//----------------- \nNODE:\nConnections:")
    n.printConnections()
    fmt.Println("Blockchain:")
    n.printBlockchain()
    fmt.Printf("Your Address:\n %v \n", n.address)
    fmt.Printf("Seed Adddress:\n %v \n", n.seed)
    fmt.Println("Seen Blocks:")
    n.printSeenTrans()
    fmt.Println("-----------------//")
}

func (n Node) printSeenTrans(){
    for blockHashString, _  := range n.seenBlocks{
        blockHashBytes := []byte(blockHashString)
        fmt.Printf("  %v\n", blockHashBytes)
    }
}

func (n Node) printBlockchain(){
    for i := range n.blockchain.Blocks {
        block := n.blockchain.Blocks[i]
        fmt.Printf(" Block %d is: \n  PrevHash: %v \n  Info:     %v \n  Hash:     %v \n", i, block.PrevHash, block.Info, block.Hash)
    }
}

func (n Node) printConnections(){
    for conn, id := range n.connections {
        localAddr := conn.LocalAddr().String()
        remoteAddr := conn.RemoteAddr().String()
        fmt.Printf(" ID: %v, Connection: %v to %v \n", id, localAddr, remoteAddr)
    }
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

func newNode() Node {
    myNode := Node{make(map[net.Conn]int), 0, Blockchain{[]Block{genesisBlock}}, "", "", map[string]bool{}}
    return myNode
}

func (myNode Node) run(listenPort string, seedPort string, joinFlag bool) {
    
    // specify ports to seed and listen to
    myNode.updateAddress(listenPort)
    myNode.updateSeed(seedPort)

    // create channels
    inputChannel            := make(chan string)
    transmissionChannel     := make(chan *Transmission)
    connChannel             := make(chan net.Conn)
    disconnChannel          := make(chan net.Conn)
    requestChannel          := make(chan net.Conn)
    addressesChannel        := make(chan []string)
    blockChannel            := make(chan Block)
    blockchainRequestChannel:= make(chan net.Conn)
    blockchainChannel       := make(chan Blockchain)

    startListening(listenPort, connChannel, inputChannel)
    if joinFlag { // if the user requested to join a seed node // need to make sure you can't join if you don't supply a seed
        fmt.Println("Dialing seed node at port " + seedPort + "...")
         go dialNode(seedPort, connChannel)
    }

    // handle go routines
    for {
        select {
            case conn    := <- connChannel: // listener picked up new conn
                myNode.incrementConnID()
                myNode.connections[conn] = myNode.nextConnID // assign connection an ID
                go listenToConn(conn, transmissionChannel, disconnChannel, requestChannel, addressesChannel, blockchainRequestChannel, blockchainChannel)

            case disconn := <- disconnChannel: // established connection disconnected
                connID := myNode.connections[disconn]
                delete(myNode.connections, disconn) // remove the connection from the nodes list of connections
                fmt.Printf("* Connection %v has been disconnected \n", connID)

            case trans := <- transmissionChannel:  // new transmission sent to node // handles adding, validating, and sending blocks to network
                notMinedAndValid   := myNode.seenBlocks[string(trans.Block.Hash)] == false  && trans.BeenSent == true && myNode.blockchain.isValidBlock(trans.Block)
                notMinedAndInvalid := myNode.seenBlocks[string(trans.Block.Hash)] == false  && trans.BeenSent == true && !myNode.blockchain.isValidBlock(trans.Block)
                minedButNotSent    := myNode.seenBlocks[string(trans.Block.Hash)] == true   && trans.BeenSent == false
                alreadySent        := myNode.seenBlocks[string(trans.Block.Hash)] == true   && trans.BeenSent == true
                if notMinedAndValid {
                    myNode.seenBlocks[string(trans.Block.Hash)] = true
                    myNode.blockchain.addBlock(trans.Block)
                    fmt.Printf("[notMinedAndValid] Added block #%v sent from network to my blockchain, and sending it to network\n", trans.Block.Index)
                    trans.updateSender(myNode.address)
                    forwardTransToNetwork(*trans, myNode.connections) // forward messages to the rest of network
                } else if notMinedAndInvalid {
                    myNode.seenBlocks[string(trans.Block.Hash)] = true
                    myBlockchainLength := myNode.blockchain.getLastBlock().Index
                    if trans.Block.Index > myBlockchainLength {
                        connThatSentHigherBlockIndex := myNode.getConnForAddress(trans.Sender)
                        fmt.Println("I was sent a block with a higher index, now requesting full chain to validate")
                        requestBlockchain(connThatSentHigherBlockIndex)
                    }
                    fmt.Printf("[notMinedAndInvalid] Did not add block #%v sent from network to my chain, did not forward\n", trans.Block.Index)
                } else if minedButNotSent { //mined but not sent out yet,
                    trans.updateBeenSent()
                    trans.updateSender(myNode.address) 
                    fmt.Printf("[minedButNotSent] Mined block #%v, sending to network\n", trans.Block.Index)
                    forwardTransToNetwork(*trans, myNode.connections) // forward messages to the rest of network
                } else if alreadySent{
                    fmt.Printf("[alreadySent] Already seen block #%v, did not forward", trans.Block.Index)
                } else {
                    fmt.Println("Some other case, this should not occur:")
                }

            case conn := <-  requestChannel:  // was requested addresses to send
                addressesToSendTo := myNode.getRemoteAddresses()
                sendConnectionsToNode(conn, addressesToSendTo)

            case addresses := <- addressesChannel:  //received addresses to add
                fmt.Print("Seed node sent these addresses to connect to:\n-> " )
                fmt.Println(addresses)
                approvedAddresses := []string{}
                for i := range addresses {
                    r, _ := regexp.Compile(":.*") // match everything after the colon
                    port := r.FindString(addresses[i])
                    if len(port) == 5 {  // in a real network this should simply be 1999
                        go dialNode(port, connChannel)
                        approvedAddresses = append(approvedAddresses, addresses[i])
                    }
                }
                fmt.Print("These new connections will be added:\n->")
                fmt.Println(approvedAddresses)

            case conn    := <- blockchainRequestChannel:
                sendBlockchainToNode(conn, myNode.blockchain)

            case blockchain := <- blockchainChannel: // node was sent a blockchain
                fmt.Println("You were sent a blockchain")
                if blockchain.isValidChain() {
                    myNode.blockchain = blockchain
                    fmt.Println("Blockchain accepted: ")
                    fmt.Println(blockchain)
                } else {
                    fmt.Println("Blockchain rejected, invalid")
                }

            case block   := <- blockChannel: // new block was mined (only mined blocks sent here)
                if myNode.blockchain.isValidBlock(block){
                    myNode.blockchain.addBlock(block)
                    myNode.seenBlocks[string(block.Hash)] = true // specify weve now seen this block but don't update the trans address until its processed there
                    go sendTransFromMinedBlock(block, transmissionChannel)
                } else {
                    fmt.Printf("Did not add mined block #%v\n", block.Index)
                }
                go myNode.blockchain.mineBlock(blockChannel)

            case input   := <- inputChannel: // user entered some input
                outgoingArgs := strings.Fields(strings.Split(input,"\n")[0]) // remove newline char and seperate into array by whitespace
                arg0 := strings.ToLower(outgoingArgs[0])
                switch arg0 {
                case "mine":
                    go myNode.blockchain.mineBlock(blockChannel)                        
                case "getchain":
                    if myNode.seed == "" {
                        fmt.Println("You must have a seed node to request a blockchain")
                    } else{
                        seedConn := myNode.getConnForAddress(myNode.seed)
                        requestBlockchain(seedConn)                        
                    }
                case "getconns":
                    if myNode.hasConnectionOfAddress(myNode.seed){
                        seedConn := myNode.getConnForAddress(myNode.seed)
                        fmt.Println("Requesting more connections from seed " + myNode.seed + " ...")
                        requestConnections(seedConn)
                    } else {
                        fmt.Println("You are not connected to your seed node to make a request..")
                    }
                case "node":
                    myNode.printNode()
                case "help":
                    showNodeHelp()
                default:
                    fmt.Println("Enter 'help' for options.")
                }
        }

    }
}

/*----------------------*
 * FUNCTION DEFINITIONS *
 *----------------------*/

func startListening (port string, connChannel chan net.Conn, inputChannel chan string) {
    listener, err := net.Listen("tcp", port)
    if err != nil {
        fmt.Println("There was an error setting up the listener:")
        fmt.Println(err)
    }

    go acceptConn(listener, connChannel)
    go listenForUserInput(inputChannel)
}

func acceptConn (listener net.Listener, connChannel chan net.Conn) {
    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Printf("There was an error accepting a new connection: \n %v", err)
        }
        fmt.Println("* Listener accepted connection through port " + conn.LocalAddr().String() + " from " + conn.RemoteAddr().String())
        connChannel <- conn //send to conection channel
    }
}

func dialNode (port string, connChannel chan net.Conn) {
    conn, err := net.Dial("tcp", getPrivateIP()+port)
    if err != nil {
        fmt.Println("**Make sure there is someone listening at " + port + "**")
        fmt.Println(err)
    }
    fmt.Println("Connetion established out of port " + conn.LocalAddr().String() + " dialing to " + conn.RemoteAddr().String())
    connChannel <- conn
}

func listenForUserInput (inputChannel chan string) {
    for {
        reader := bufio.NewReader(os.Stdin) //constantly be reading in from std in
        input, err := reader.ReadString('\n')
        if (err != nil || input == "\n") {
            // fmt.Println(err)
        } else {
            fmt.Println()
            inputChannel <- input
        }
    }
}

func listenToConn (conn                          net.Conn, 
                   transmissionChannel      chan *Transmission,
                   disconnChannel           chan net.Conn,
                   requestChannel           chan net.Conn,
                   addressesChannel         chan []string,
                   blockchainRequestChannel chan net.Conn,
                   blockchainChannel        chan Blockchain) {
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
            transmissionChannel <- &communication.Trans
        case 1:
            addressesChannel <- communication.SentAddresses
        case 2:
            fmt.Println("You have been requested to send your connection addresses to a peer at " + conn.RemoteAddr().String() + " ...")
            requestChannel <- conn
        case 3:
            blockchainChannel <- communication.Blockchain
        case 4:
            fmt.Println("You have been requested to send your blockchain address to a peer at " + conn.RemoteAddr().String() + " ...")
            blockchainRequestChannel <- conn
        default:
            fmt.Println("There was a problem decoding the message")
        }
    }
    disconnChannel <- conn // disconnect must have occurred if we exit the for loop
}

func forwardTransToNetwork (trans Transmission, connections map[net.Conn]int) {
    // fmt.Println("IN SEND TRANS TO NET, trans.BeenSent is:")
    for conn, _ := range connections { // loop through all this nodes connections
        // destinationAddr := conn.RemoteAddr().String() // get the destination of the connection
        communication := Communication{0, trans, []string{}, Blockchain{}}
        encoder       := gob.NewEncoder(conn)
        encoder.Encode(communication)        
    }
}

func requestConnections (conn net.Conn){
    communication := Communication{2, Transmission{}, []string{}, Blockchain{}}
    encoder       := gob.NewEncoder(conn)
    encoder.Encode(communication)
}

func sendConnectionsToNode (conn net.Conn, addresses []string){
    communication := Communication{1, Transmission{}, addresses, Blockchain{}}
    encoder       := gob.NewEncoder(conn)
    encoder.Encode(communication)
}

func sendBlockchainToNode (conn net.Conn, blockchain Blockchain){
    communication := Communication{3, Transmission{}, []string{}, blockchain}
    encoder   := gob.NewEncoder(conn)
    encoder.Encode(communication)
    fmt.Printf("Sent my copy of blockchain to %v", conn.RemoteAddr().String())
}

func requestBlockchain (conn net.Conn){
    communication := Communication{4, Transmission{}, []string{}, Blockchain{}}
    encoder   := gob.NewEncoder(conn)
    encoder.Encode(communication)
}

func sendTransFromMinedBlock(block Block, transmissionChannel chan *Transmission){
    trans := Transmission{block, false, ""}
    transmissionChannel <- &trans
}

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



