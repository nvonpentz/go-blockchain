package main 

import (
	"testing"
	"net"
	"fmt"
)

func createTestListener(port string) *net.Listener {
	listener, err := net.Listen("tcp", port)
    if err != nil {
        fmt.Println("There was an error setting up the listener:")
        fmt.Println(err)
    }
    return &listener
}

func testDial(address string) net.Conn {
	conn, err := net.Dial("tcp", address)
    if err != nil {
        fmt.Println("**Make sure there is someone listening at " + address + "**")
        fmt.Println(err)
    }

    return conn
}

func TestIncrementConnID(t *testing.T) {
	n := newNode()
	n.incrementConnID()
	if n.nextConnID != 1 {
		t.Error("Expected 1, got %v", n.nextConnID)
	}
}

func TestStartListening(t *testing.T){
	listenPort   := ":1999"
	connChannel  := make(chan net.Conn)
	inputChannel := make(chan string)

	startListening(listenPort, connChannel, inputChannel)
	conn, err := net.Dial("tcp", "127.0.0.1" + listenPort)
	if err != nil {
		t.Error("Unable to make a connection using startListening()")
		fmt.Println(conn)
	}
	conn.Close()
}

func TestAcceptConn(t *testing.T){
	listener, err := net.Listen("tcp", ":2000")
    if err != nil {
        fmt.Println("There was an error setting up the listener:")
        fmt.Println(err)
    }
	// listener := *createTestListener(":1999")
    connChannel := make(chan net.Conn)
    go acceptConn(listener, connChannel)

    conn := testDial(":2000")
    receivedConn := <- connChannel
    if conn.LocalAddr().String() != receivedConn.RemoteAddr().String() {
    	t.Error("Could not establish connection using acceptConn().  sending and receiving conns should be the same:")
    }
	conn.Close()
}












