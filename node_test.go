package main 

import ("testing"
		"net"
		)

func TestIncrementConnID(t *testing.T) {
	n := Node{map[net.Conn]int{}, 0, Blockchain{}, "", ""}
	n.incrementConnID()
	if n.nextConnID != 1 {
		t.Error("Expected 1, got %v", n.nextConnID)
	}
}