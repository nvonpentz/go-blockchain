package main 

import(
	"testing"

	"github.com/nvonpentz/go-hashable-keys"
)

func TestVerifyPacketSignature(t *testing.T){
	keys := hashkeys.GenerateNewKeypair()
	packet := createPacket("document.txt", *keys)

	if !verifyPacketSignature(packet) {
		t.Error("fails to approve valid packet")
	}
}
