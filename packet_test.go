package main 

import(
	"testing"
)

func TestVerifyPacketSignature(t *testing.T){
	keys := hashkeys.GenerateNewKeyPair()
	packet := createPacket("document.txt", keys)

	if !verifyPacketSignature(packet) {
		t.Error("fails to approve valid packet")
	}
}
