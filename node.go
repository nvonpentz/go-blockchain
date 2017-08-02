package main

import(
    "fmt"
    "net"
    "bufio"
    "os"
    "encoding/gob"
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
}

/*
A Transmission is simply a message that will be sent throughout the network.  
It includes the actual message as well as the addresses of nodes who have
already received the Transmission
*/
type Transmission struct {
    Block Block
    VisitedAddresses map[string]bool // map for efficiency
}

/*
Every communication between nodes over the TCP connection will be 
through sending this communication object.  The ID represents the type:

0 - means we will be receiving a transmission
1 - means we will be receiving a slice of sent addresses
2 - means we were requested to send conections
3 - means we were requested to send your blockchain 
4 - means we will be receiving a blockchain
5 - means we will be requesting
*/
type Communication struct {
    ID int
    Trans Transmission
    SentAddresses []string
    Blockchain Blockchain
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

func (n *Node) replaceChainIfLonger(theirChain Blockchain){
    myLastBlock    := n.blockchain.getLastBlock()
    theirLastBlock := theirChain.getLastBlock()
    if myLastBlock.Index < theirLastBlock.Index && theirChain.isValidChain() {
        n.blockchain = theirChain
        fmt.Printf("Replaced my blockchain with")
    } else {
        fmt.Printf("Did not replace blockchain")
    }
}

func (n Node) printBlockchain(){
    for i := range n.blockchain.Blocks {
        block := n.blockchain.Blocks[i]
        fmt.Printf(" Block %d is: \n  PrevHash: %v \n  Info:     %v \n  Hash:     %v \n", i, block.PrevHash, block.Info, block.Hash)
    }
}

func (n Node) listConnections(){
    for conn, id := range n.connections {
        localAddr := conn.LocalAddr().String()
        remoteAddr := conn.RemoteAddr().String()
        fmt.Printf(" ID: %v, Connection: %v to %v \n", id, localAddr, remoteAddr)
    }
}

func (n Node) numberConns() (int){
    return len(n.connections)
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

func (t *Transmission) updateVisitedAddresses(address string) {
    t.VisitedAddresses[address] = true
}

func (t *Transmission) hasAddress(address string) bool {
    if val := !t.VisitedAddresses[address]; val {
        return false
    }
    return true
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
            fmt.Println("You have been requested to send your conenction address to a peer at " + conn.RemoteAddr().String() + " ...")
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
    disconnChannel <- conn
}

func forwardTransToNetwork (trans Transmission, connections map[net.Conn]int) {
    for conn, _ := range connections { // loop through all this nodes connections
        destinationAddr := conn.RemoteAddr().String() // get the destination of the connection
        if !trans.hasAddress(destinationAddr){ // look to see if this transmission has already been to this node so we don't send to any we KNOW have seen it.
            communication := Communication{0, trans, []string{}, Blockchain{}}
            encoder       := gob.NewEncoder(conn)
            encoder.Encode(communication)        
            fmt.Printf("Sent transsmission to network:\n %v to %v \n", trans, destinationAddr)
        }
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
    fmt.Println("requested blockchain")
}

func sendTransFromMinedBlock(block Block, transmissionChannel chan *Transmission){
    trans := Transmission{block, map[string]bool{}}
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
