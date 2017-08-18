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
    blockchain  Blockchain
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
    blockWrapperChannel      := make(chan *BlockWrapper)
    newConnChannel           := make(chan net.Conn) // new connections added
    disconChannel            := make(chan net.Conn) // new disconnestion
    connRequestChannel       := make(chan net.Conn) // received a request to send connections 
    sentAddressesChannel     := make(chan []string) // received addresses to make connections
    minedBlockChannel        := make(chan Block)    // new block was mined
    blockchainRequestChannel := make(chan net.Conn)
    sentBlockchainChannel    := make(chan Blockchain)

    // listen to user input
    go listenForUserInput(minedBlockChannel, blockWrapperChannel, &myNode)

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
                go listenToConn(conn, blockWrapperChannel, disconChannel, connRequestChannel, sentAddressesChannel, blockchainRequestChannel, sentBlockchainChannel)

            case discon       := <- disconChannel: // established connection disconnected
                connID := myNode.connections[discon]
                delete(myNode.connections, discon) // remove the connection from the nodes list of connections
                fmt.Printf("* Connection %v has been disconnected \n", connID)

            case blockWrapper := <- blockWrapperChannel:  // new blockWrapper sent to node // handles adding, validating, and sending blocks to network
                block  := blockWrapper.Block
                // fmt.Printf("here is seen blocks:\n %v\n", myNode.seenBlocks)
                // fmt.Printf("here is the hash of the newly mined block: %v\n", block.Hash)
                seenBlock := myNode.seenBlocks[string(block.Hash)] == true
                // fmt.Printf("here is seenBock %v\n", seenBlock)
                if !seenBlock {
                    blockValid := myNode.blockchain.isValidBlock(block)
                    if blockValid {
                        myNode.seenBlocks[string(block.Hash)] = true // only set to seen if we validate it, otherwise it will come around again
                        myNode.forwardBlockWrapperToNetwork(BlockWrapper{Block: block, Sender: myNode.address}, myNode.connections)                        
                        myNode.blockchain.addBlock(block)
                        fmt.Println("sent blockchain to network")
                    } else {
                        fmt.Println("block was not considered valid, making request for whole chain to compare..")                        
                        requestBlockchain(myNode.getConnForAddress(blockWrapper.Sender)) //request blockchain ending in block, ba                            
                    }
                } else {
                    fmt.Println("seen this block before, ignoring..")
                }
                // myNode.handleBlockWrapper(blockWrapper)
            case conn         := <-  connRequestChannel:  // was requested addresses to send
                addressesToSendTo := myNode.getRemoteAddresses()
                sendConnectionsToNode(conn, addressesToSendTo)

            case addresses    := <- sentAddressesChannel:  //received addresses to add
                fmt.Printf("Seed node sent these addresses to connect to:\n-> %v\n", addresses)
                myNode.handleSentAddresses(addresses, newConnChannel)

            case conn         := <- blockchainRequestChannel:
                sendBlockchainToNode(conn, myNode.blockchain)

            case blockchain   := <- sentBlockchainChannel: // node was sent a blockchain
                myNode.handleSentBlockchain(blockchain, blockWrapperChannel)
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
func (n *Node) forwardBlockWrapperToNetwork(blockWrapper BlockWrapper, connections map[net.Conn]int) {
    for conn, _ := range connections { // loop through all this nodes connections
        // destinationAddr := conn.RemoteAddr().String() // get the destination of the connection
        communication := Communication{0, blockWrapper, []string{}, Blockchain{}}
        encoder       := gob.NewEncoder(conn)
        encoder.Encode(communication)        
    }
}

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

func (n *Node) handleSentBlockchain(blockchain Blockchain, blockWrapperChannel chan *BlockWrapper){
    fmt.Println("You were sent a blockchain")
    if blockchain.isValidChain() {
        lastIndex := len(blockchain.Blocks)-1
        semiReplacementChain := Blockchain{blockchain.Blocks[:lastIndex]}
        n.blockchain = semiReplacementChain
        
        seenBlocks := make(map[string]bool)  // need a new set of seen blocks associated with 
        for _ , b := range semiReplacementChain.Blocks{
            seenBlocks[string(b.Hash)] = true
        }
        n.seenBlocks = seenBlocks //replace with the associated seen blocks

        fmt.Println("Blockchain accepted: ")
        fmt.Println(blockchain)

        lastBlock := blockchain.Blocks[lastIndex]
        blockWrapper := BlockWrapper{Block: lastBlock, Sender: n.address}
        go func () {blockWrapperChannel <- &blockWrapper}()
        fmt.Println("sent the tip of the replacement chain to the blockchannel")
    } else {
        fmt.Println("Blockchain rejected, invalid!")
    }
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

func (n Node) printSeenBlockWrapper(){
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
    n.printSeenBlockWrapper()
    fmt.Println(" Blockchain:")
    n.printBlockchain()
    fmt.Println("*------------------*")
}

/*-------------*
 *  FUNCTIONS  * 
 *-------------*/

func newNode() Node {
    myNode := Node{make(map[net.Conn]int), 0, Blockchain{[]Block{genesisBlock}}, "", "", map[string]bool{}}
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
                            blockWrapperChannel      chan *BlockWrapper,
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
            blockWrapperChannel <- &communication.BlockWrapper
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
            break
        }
    }
    disconChannel <- conn // disconnect must have occurred if we exit the for loop
}

func requestConnections(conn net.Conn){
    communication := Communication{2, BlockWrapper{}, []string{}, Blockchain{}}
    encoder       := gob.NewEncoder(conn)
    encoder.Encode(communication)
}

func sendConnectionsToNode(conn net.Conn, addresses []string){
    communication := Communication{1, BlockWrapper{}, addresses, Blockchain{}}
    encoder       := gob.NewEncoder(conn)
    encoder.Encode(communication)
}

func sendBlockchainToNode(conn net.Conn, blockchain Blockchain){
    communication := Communication{3, BlockWrapper{}, []string{}, blockchain}
    encoder   := gob.NewEncoder(conn)
    encoder.Encode(communication)
    fmt.Printf("Sent my copy of blockchain to %v", conn.RemoteAddr().String())
}

func requestBlockchain(conn net.Conn){
    communication := Communication{4, BlockWrapper{}, []string{}, Blockchain{}}
    encoder   := gob.NewEncoder(conn)
    encoder.Encode(communication)
}

func sendBlockWrapperFromMinedBlock(block Block, blockWrapperChannel chan *BlockWrapper){
    blockWrapper := BlockWrapper{block, ""}
    blockWrapperChannel <- &blockWrapper
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

