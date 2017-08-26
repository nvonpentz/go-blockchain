package main 

import(
	"testing"

	"github.com/nvonpentz/go-hashable-keys"
)

func TestVerifyPacketSignature(t *testing.T){
	keys01 := hashkeys.GenerateNewKeypair()
	keys02 := hashkeys.GenerateNewKeypair()

	packet01 := createPacket("document.txt", *keys01)

	// valid packet
	if !verifyPacketSignature(packet01) {
		t.Error("fails to approve valid packet")
	}


	// invalid packet; wrong public key for signature
	doc       := readDocument("document.txt")
	hashedDoc := hashDocument(doc)
	signature := signHash(hashedDoc, *keys01) // sign with keys01

	packet02 := Packet{Hash: hashedDoc, Signature: signature, Owner: keys02.Public}

	if verifyPacketSignature(packet02) {
		t.Error("validates document with wrong public key for signature")
	}

}

func TestVerifyPacketList(t *testing.T){
	// base case
	if !verifyPacketList([]Packet{}){
		t.Error("fails to validate empty packet list")
	}

	// create 3 valid packets
	keys01 := hashkeys.GenerateNewKeypair()
	keys02 := hashkeys.GenerateNewKeypair()
	keys03 := hashkeys.GenerateNewKeypair()

	packet01 := createPacket("document.txt", *keys01)
	packet02 := createPacket("document.txt", *keys02)
	packet03 := createPacket("document.txt", *keys03)

	packets := []Packet{packet01, packet02, packet03}
	if !verifyPacketList(packets){
		t.Error("fail to validate valid packet list of length 3")
	}

	// create an invalid packet
	doc       := readDocument("document.txt")
	hashedDoc := hashDocument(doc)
	signature := signHash(hashedDoc, *keys01) // sign with keys01

	packet04 := Packet{Hash: hashedDoc, Signature: signature, Owner: keys02.Public}

	packets = []Packet{packet01, packet04, packet03}

	if verifyPacketList(packets){
		t.Error("validates invalid packet list")
	}
}









