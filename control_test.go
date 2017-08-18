package main 

import(
	"testing"
	"net"
	"fmt"
	"time"
)

func TestListenForConnections(t *testing.T){
	listenPort       := ":2000" //specific 
	newConnChannel   := make(chan net.Conn)

	listenForConnections(listenPort, newConnChannel)
	conn, err := net.Dial("tcp", "127.0.0.1" + listenPort)
	if err != nil {
		t.Error("Unable to make a connection using listenForUserInput()")
		fmt.Println(conn)
	}
	conn.SetDeadline(time.Now())
}