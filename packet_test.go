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

func TestEqualPackets(t *testing.T){
	// create 3 valid packets
	keys01 := hashkeys.GenerateNewKeypair()
	keys02 := hashkeys.GenerateNewKeypair()

	packet01 := createPacket("document.txt", *keys01)
	packet02 := createPacket("document.txt", *keys02)

	if !equalPackets(packet01, packet01){
		t.Error("fails to validate equal packets")
	}

	if equalPackets(packet01, packet02){
		t.Error("validates to different packets")
	}
}

func TestPacketListHasPacket(t *testing.T){
	// create 4 different packets
	keys01 := hashkeys.GenerateNewKeypair()
	keys02 := hashkeys.GenerateNewKeypair()
	keys03 := hashkeys.GenerateNewKeypair()
	keys04 := hashkeys.GenerateNewKeypair()


	packet01 := createPacket("document.txt", *keys01)
	packet02 := createPacket("document.txt", *keys02)
	packet03 := createPacket("document.txt", *keys03)
	packet04 := createPacket("document.txt", *keys04)

	packets := []Packet{packet01, packet02, packet03}

	if !packetListHasPacket(packets, packet01){
		t.Error("fails to find packet in list")
	}
	if !packetListHasPacket(packets, packet03){
		t.Error("fails to find packet in list")
	}
	if packetListHasPacket(packets, packet04){
		t.Error("claims to find packet not in list")
	}
}

func TestPacketListHasPacketHash(t *testing.T){
	// TO DO
}

func TestGetPackFromListByHash(t *testing.T){
	// TO DO
}


