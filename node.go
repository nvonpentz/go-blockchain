package main

import(
    "fmt"
    "net"
    "bufio"
    "os"
    "encoding/gob"
    "strings"
    "regexp"
    "net/http"
    "io/ioutil"
)

type Node struct {
    connections map[net.Conn]int
    nextConnID int
    blockchain Blockchain
    address string
    seed string
    seenBlocks map[string]bool
}

func newNode() Node {
    myNode := Node{make(map[net.Conn]int), 0, Blockchain{[]Block{genesisBlock}}, "", "", map[string]bool{}}
    return myNode
}

func (n *Node) updatePorts(listenPort string, seedInfo string, publicFlag bool) {
    if publicFlag{
        n.seed = seedInfo + ":1999" // if public ip, seed is specifiec seedInfo:1999
        n.address = n.getPublicIP() + ":1999" // must set up port forwarding
    } else { 
        n.seed = n.getPrivateIP() + ":" + seedInfo  // no default seed
        n.address = n.getPrivateIP() + listenPort
    }
}

func (n *Node) startListening(port string, newConnChannel chan net.Conn, userInputChannel chan string) {
    listener, err := net.Listen("tcp", port)
    if err != nil {
        fmt.Println("There was an error setting up the listener:")
        fmt.Println(err)
    }

    go n.acceptConn(listener, newConnChannel)
    go n.listenForUserInput(userInputChannel)
}

func (n *Node) acceptConn(listener net.Listener, newConnChannel chan net.Conn) {
    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Printf("There was an error accepting a new connection: \n %v", err)
        }
        fmt.Println("* Listener accepted connection through port " + conn.LocalAddr().String() + " from " + conn.RemoteAddr().String())
        newConnChannel <- conn //send to conection channel
    }
}

func (n *Node) dialNode(address string, newConnChannel chan net.Conn) {
    conn, err := net.Dial("tcp", address)
    if err != nil {
        fmt.Println("**Make sure there is someone listening at " + address + "**")
        fmt.Println(err)
    }
    fmt.Println("Connection established out of port " + conn.LocalAddr().String() + " dialing to " + conn.RemoteAddr().String())
    newConnChannel <- conn
}

func (n *Node) listenForUserInput(userInputChannel chan string) {
    for {
        reader := bufio.NewReader(os.Stdin) //constantly be reading in from std in
        input, err := reader.ReadString('\n')
        if (err != nil || input == "\n") {
        } else {
            fmt.Println()
            userInputChannel <- input
        }
    }
}

func (n *Node) listenToConn(conn                          net.Conn, 
                            transmissionChannel      chan *Transmission,
                            disconChannel            chan net.Conn,
                            connRequestChannel       chan net.Conn,
                            sentAddressesChannel     chan []string,
                            blockchainRequestChannel chan net.Conn,
                            sentBlockchainChannel    chan Blockchain) {
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
            sentAddressesChannel <- communication.SentAddresses
        case 2:
            fmt.Println("You have been requested to send your connection addresses to a peer at " + conn.RemoteAddr().String() + " ...")
            connRequestChannel <- conn
        case 3:
            sentBlockchainChannel <- communication.Blockchain
        case 4:
            fmt.Println("You have been requested to send your blockchain address to a peer at " + conn.RemoteAddr().String() + " ...")
            blockchainRequestChannel <- conn
        default:
            fmt.Println("There was a problem decoding the message")
        }
    }
    disconChannel <- conn // disconnect must have occurred if we exit the for loop
}

func (n *Node) forwardTransToNetwork(trans Transmission, connections map[net.Conn]int) {
    for conn, _ := range connections { // loop through all this nodes connections
        // destinationAddr := conn.RemoteAddr().String() // get the destination of the connection
        communication := Communication{0, trans, []string{}, Blockchain{}}
        encoder       := gob.NewEncoder(conn)
        encoder.Encode(communication)        
    }
}

func (n *Node) handleTrans(trans *Transmission){
    
    notMinedAndValid   := n.seenBlocks[string(trans.Block.Hash)] == false  && trans.BeenSent == true && n.blockchain.isValidBlock(trans.Block)
    notMinedAndInvalid := n.seenBlocks[string(trans.Block.Hash)] == false  && trans.BeenSent == true && !n.blockchain.isValidBlock(trans.Block)
    minedButNotSent    := n.seenBlocks[string(trans.Block.Hash)] == true   && trans.BeenSent == false
    alreadySent        := n.seenBlocks[string(trans.Block.Hash)] == true   && trans.BeenSent == true
    
    if notMinedAndValid {
        n.seenBlocks[string(trans.Block.Hash)] = true
        n.blockchain.addBlock(trans.Block)
        fmt.Printf("[notMinedAndValid] Added block #%v sent from network to my blockchain, and sending it to network\n", trans.Block.Index)
        trans.updateSender(n.address)
        n.forwardTransToNetwork(*trans, n.connections) // forward messages to the rest of network
    } else if notMinedAndInvalid {
        n.seenBlocks[string(trans.Block.Hash)] = true
        myBlockchainLength := n.blockchain.getLastBlock().Index
        if trans.Block.Index > myBlockchainLength {
            connThatSentHigherBlockIndex := n.getConnForAddress(trans.Sender)
            fmt.Println("I was sent a block with a higher index, now requesting full chain to validate")
            n.requestBlockchain(connThatSentHigherBlockIndex)
        }
        fmt.Printf("[notMinedAndInvalid] Did not add block #%v sent from network to my chain, did not forward\n", trans.Block.Index)
    } else if minedButNotSent { //mined but not sent out yet,
        trans.updateBeenSent()
        trans.updateSender(n.address) 
        fmt.Printf("[minedButNotSent] Mined block #%v, sending to network\n", trans.Block.Index)
        n.forwardTransToNetwork(*trans, n.connections) // forward messages to the rest of network
    } else if alreadySent{
        fmt.Printf("[alreadySent] Already seen block #%v, did not forward\n", trans.Block.Index)
    } else {
        fmt.Println("Some other case, this should not occur:")
    }
}

func (n *Node) handleUserInput(input string, minedBlockChannel chan Block) {
    outgoingArgs := strings.Fields(strings.Split(input,"\n")[0]) // remove newline char and seperate into array by whitespace
    arg0 := strings.ToLower(outgoingArgs[0])
    switch arg0 {
    case "mine":
        go n.blockchain.mineBlock(minedBlockChannel)                        
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

func (n *Node) handleSentAddresses(addresses []string, newConnChannel chan net.Conn){
    approvedAddresses := []string{}
    for i := range addresses {
        r, _ := regexp.Compile(":.*") // match everything after the colon
        port := r.FindString(addresses[i])
        if len(port) == 5 {  // in a real network this should simply be 1999
            go n.dialNode(addresses[i], newConnChannel)
            approvedAddresses = append(approvedAddresses, addresses[i])
        }
    }
    fmt.Printf("These new connections will be added:\n->%v\n", approvedAddresses)
}

func (n *Node) handleMinedBlock(block Block, minedBlockChannel chan Block, transmissionChannel chan *Transmission) {
    if n.blockchain.isValidBlock(block){
        n.blockchain.addBlock(block)
        n.seenBlocks[string(block.Hash)] = true // specify weve now seen this block but don't update the trans address until its processed there
        go n.sendTransFromMinedBlock(block, transmissionChannel)
    } else {
        fmt.Printf("Did not add mined block #%v\n", block.Index)
    }
    go n.blockchain.mineBlock(minedBlockChannel)
}

func (n *Node) handleSentBlockchain(blockchain Blockchain){
    fmt.Println("You were sent a blockchain")
    if blockchain.isValidChain() {
        n.blockchain = blockchain
        fmt.Println("Blockchain accepted: ")
        fmt.Println(blockchain)
    } else {
        fmt.Println("Blockchain rejected, invalid")
    }
}

func (n *Node) requestConnections(conn net.Conn){
    communication := Communication{2, Transmission{}, []string{}, Blockchain{}}
    encoder       := gob.NewEncoder(conn)
    encoder.Encode(communication)
}

func (n *Node) sendConnectionsToNode(conn net.Conn, addresses []string){
    communication := Communication{1, Transmission{}, addresses, Blockchain{}}
    encoder       := gob.NewEncoder(conn)
    encoder.Encode(communication)
}

func (n *Node) sendBlockchainToNode(conn net.Conn, blockchain Blockchain){
    communication := Communication{3, Transmission{}, []string{}, blockchain}
    encoder   := gob.NewEncoder(conn)
    encoder.Encode(communication)
    fmt.Printf("Sent my copy of blockchain to %v", conn.RemoteAddr().String())
}

func (n *Node) requestBlockchain(conn net.Conn){
    communication := Communication{4, Transmission{}, []string{}, Blockchain{}}
    encoder   := gob.NewEncoder(conn)
    encoder.Encode(communication)
}

func (n *Node) sendTransFromMinedBlock(block Block, transmissionChannel chan *Transmission){
    trans := Transmission{block, false, ""}
    transmissionChannel <- &trans
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

func (n *Node) getPrivateIP() string {
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

func (n *Node) getPublicIP() string {
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

func (n Node) printSeenTrans(){
    for blockHashString, _  := range n.seenBlocks{
        blockHashBytes := []byte(blockHashString)
        fmt.Printf("  %v\n", blockHashBytes)
    }
}

func (n Node) printBlockchain(){
    for i := range n.blockchain.Blocks {
        block := n.blockchain.Blocks[i]
        fmt.Printf("  Block %d is: \n   PrevHash: %v \n   Info:     %v \n   Hash:     %v \n", i, block.PrevHash, block.Info, block.Hash)
    }
}

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
    n.printSeenTrans()
    fmt.Println(" Blockchain:")
    n.printBlockchain()
    fmt.Println("*------------------*")
}

func (myNode Node) run(listenPort string, seedInfo string, publicFlag bool) {
    joinFlag := false
    if seedInfo != "" { joinFlag = true }
    
    // specify ports to seed and listen to
    myNode.updatePorts(listenPort, seedInfo, publicFlag)

    // create channels
    userInputChannel         := make(chan string)
    transmissionChannel      := make(chan *Transmission)
    newConnChannel           := make(chan net.Conn) // new connections added
    disconChannel            := make(chan net.Conn) // new disconnestion
    connRequestChannel       := make(chan net.Conn) // received a request to send connections 
    sentAddressesChannel     := make(chan []string) // received addresses to make connections
    minedBlockChannel        := make(chan Block)    // new block was mined
    blockchainRequestChannel := make(chan net.Conn)
    sentBlockchainChannel    := make(chan Blockchain)

    myNode.startListening(listenPort, newConnChannel, userInputChannel)
    if joinFlag { // if the user requested to join a seed node // need to make sure you can't join if you don't supply a seed
        fmt.Println("Dialing seed node at port " + seedInfo + "...")
         go myNode.dialNode(myNode.seed, newConnChannel)
    }

    myNode.printNode()

    // handle go routines
    for {
        select {
            case conn    := <- newConnChannel: // listener picked up new conn
                myNode.nextConnID = myNode.nextConnID + 1
                myNode.connections[conn] = myNode.nextConnID // assign connection an ID
                go myNode.listenToConn(conn, transmissionChannel, disconChannel, connRequestChannel, sentAddressesChannel, blockchainRequestChannel, sentBlockchainChannel)

            case disconn := <- disconChannel: // established connection disconnected
                connID := myNode.connections[disconn]
                delete(myNode.connections, disconn) // remove the connection from the nodes list of connections
                fmt.Printf("* Connection %v has been disconnected \n", connID)

            case trans := <- transmissionChannel:  // new transmission sent to node // handles adding, validating, and sending blocks to network
                myNode.handleTrans(trans)

            case conn := <-  connRequestChannel:  // was requested addresses to send
                addressesToSendTo := myNode.getRemoteAddresses()
                myNode.sendConnectionsToNode(conn, addressesToSendTo)

            case addresses := <- sentAddressesChannel:  //received addresses to add
                fmt.Printf("Seed node sent these addresses to connect to:\n-> %v\n", addresses)
                myNode.handleSentAddresses(addresses, newConnChannel)

            case conn    := <- blockchainRequestChannel:
                myNode.sendBlockchainToNode(conn, myNode.blockchain)

            case blockchain := <- sentBlockchainChannel: // node was sent a blockchain
                myNode.handleSentBlockchain(blockchain)

            case block   := <- minedBlockChannel: // new block was mined (only mined blocks sent here)
                myNode.handleMinedBlock(block, minedBlockChannel, transmissionChannel)

            case input   := <- userInputChannel: // user entered some input
                myNode.handleUserInput(input, minedBlockChannel)
        }

    }
}




