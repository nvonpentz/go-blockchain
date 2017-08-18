package main 

import (
	"testing"
	"net"
	"fmt"
	// "time"
	// "encoding/gob"
)

func createTestListener(port string) *net.Listener {
	listener, err := net.Listen("tcp", port)
    if err != nil {
        fmt.Println("There was an error setting up the listener:")
        fmt.Println(err)
    }
    return &listener
}


func TestAcceptConn(t *testing.T){
	listenPort       := ":1999"
	newConnChannel   := make(chan net.Conn)
	// userInputChannel := make(chan string)

	listener, err := net.Listen("tcp", listenPort)
    if err != nil {
        fmt.Println("There was an error setting up the listener:")
        fmt.Println(err)
    }
    go acceptConn(listener, newConnChannel)

	conn1, err := net.Dial("tcp", "127.0.0.1" + listenPort)
	if err != nil {
		t.Error("Unable to make a connection using acceptConn()")
		fmt.Println(conn1)
	}

	conn2, err := net.Dial("tcp", "127.0.0.1" + listenPort)
	if err != nil {
		t.Error("Unable to make a connection using acceptConn()")
		fmt.Println(conn2)
	}

	conn1.Close()
	conn2.Close()
	listener.Close()
}

func TestDialNode(t *testing.T){
	listenPort       := ":1999"
	newConnChannel   := make(chan net.Conn)

	listener, err := net.Listen("tcp", listenPort)
    if err != nil {
        fmt.Println("There was an error setting up the listener:")
        fmt.Println(err)
    }
	go dialNode("127.0.0.1:1999", newConnChannel)
	acceptedConn, err := listener.Accept()
	if err != nil {
		t.Error("Unable to make a connection using n.dialNode()")
	}
	deliveredConn := <- newConnChannel
	if deliveredConn.LocalAddr().String() != acceptedConn.RemoteAddr().String() {
		t.Error("Unable to make a connection using n.DialNode")
    }
    listener.Close()
    deliveredConn.Close()
    acceptedConn.Close()
}

// func TestListenToConn(t *testing.T){
// 	listenPort := ":1999"

// 	// create channels
//     blockWrapperChannel      := make(chan *BlockWrapper)
//     disconChannel            := make(chan net.Conn) // new disconnestion
//     connRequestChannel       := make(chan net.Conn) // received a request to send connections 
//     sentAddressesChannel     := make(chan []string) // received addresses to make connections
//     blockchainRequestChannel := make(chan net.Conn)
//     sentBlockchainChannel    := make(chan Blockchain)

//     listener, err := net.Listen("tcp", listenPort)
//     if err != nil {
//         fmt.Println("There was an error setting up the listener:")
//         fmt.Println(err)
//     }

//     connOut, err := net.Dial("tcp", "127.0.0.1" + listenPort)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
//     connIn, err := listener.Accept()

//     go listenToConn(connIn,
//     			   blockWrapperChannel,
//     			   disconChannel,
//     			   connRequestChannel,
//     			   sentAddressesChannel,
//     			   blockchainRequestChannel,
//     			   sentBlockchainChannel)

//     comm0 := newComm(0)
//     comm1 := newComm(1)
//     comm2 := newComm(2)
//     comm3 := newComm(3)
//     comm4 := newComm(4)

//     encoder := gob.NewEncoder(connIn)
//     err = encoder.Encode(&comm0)
//     if err != nil {
//         fmt.Println(err)
//         t.Error("Error receiving blockWrapper")
//     }
// 	fmt.Println("CHABA")

//     blockWrapper := <- blockWrapperChannel
//     fmt.Println(blockWrapper)
//     if blockWrapper.Sender != comm0.BlockWrapper.Sender {
//     	t.Error("blockWrapper.Sender != comm0.BlockWrapper.Sender")
//     }

//     err = encoder.Encode(&comm1)
//     if err != nil {
//         fmt.Println(err)
//         t.Error("Error receiving sent addresses")
//     }

//     err = encoder.Encode(&comm2)
//     if err != nil {
//         fmt.Println(err)
//         t.Error("Error sending connections")
//     }
//     err = encoder.Encode(&comm3)
//     if err != nil {
//         fmt.Println(err)
//         t.Error("Error receiving blockchain")
//     }
//     err = encoder.Encode(&comm4)
//     if err != nil {
//         fmt.Println(err)
//         t.Error("Error sending blockchain")
//     }
//     err = encoder.Encode("disconnect")
//     if err != nil {
//         fmt.Println(err)
//     }
//     discon := <- disconChannel
//     if discon.LocalAddr().String() != discon.RemoteAddr().String() {
//     	t.Error("disconnection and connection do not align")
//     }
// 	listener.Close()
// 	connIn.Close()
// 	connOut.Close()
// }

// func TestForwardBlockWrapperToNewtork(t *testing.T){
// 	n := newNode()
// 	listenPort       := ":1999"
// 	newConnChannel   := make(chan net.Conn)
// 	blockWrapper := emptyBlockWrapper()

// 	listener, err := net.Listen("tcp", listenPort)
//     if err != nil {
//         fmt.Println("There was an error setting up the listener:")
//         fmt.Println(err)
//     }
	
// 	go dialNode("127.0.0.1:1999", newConnChannel)
// 	conn1, err := listener.Accept()
// 	go dialNode("127.0.0.1:1999", newConnChannel)
// 	conn2, err := listener.Accept()

// 	connections := map[net.Conn]int {conn1:0, conn2:1}
// 	n.forwardBlockWrapperToNetwork(blockWrapper, connections)

//  //    var comm Communication
// 	// decoder := gob.NewDecoder(conn1)
// 	// err = decoder.Decode(&comm)
// 	// if err != nil {
// 	// 	fmt.Println(err)
// 	// }
	
// 	// fmt.Print(comm.BlockWrapper.Sender)
// 	conn1.Close()
// 	conn2.Close()
// 	listener.Close()
// }











