package main 

import(
	"crypto/sha256"
	"io/ioutil"
	"fmt"

	"github.com/nvonpentz/go-hashable-keys"
)

type Packet struct {
	Hash      []byte
	Signature []byte
	Owner     []byte
}

func readDocument(filePath string) []byte{
	document, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println(err)
	}

	return document
}

func hashDocument(document []byte) []byte{
	h := sha256.New()
	h.Write(document)
	return h.Sum(nil)
}


func signHash(hash []byte, keys hashkeys.Keypair) []byte{
	signature, err := keys.Sign(hash) //sign the hash of the transaction
	if err !=nil {
		fmt.Println(err)
	}

	return signature
}

func createPacket(filepath string, keys hashkeys.Keypair) Packet {
	document := readDocument(filepath)
	documentHash := hashDocument(document)
	signature := signHash(documentHash, keys)

	return Packet{Hash: documentHash, Signature: signature, Owner: keys.Public}
}

func verifyPacketSignature(packet Packet) bool {
	return hashkeys.SignatureVerify(packet.Owner, packet.Signature, packet.Hash)
}

func verifyPacketList(packets []Packet) bool {
	for _ , packet := range packets{
		if verifyPacketSignature(packet) == false {
			return false
		}
	}
	return true
}

func hashPacketList(list []Packet) []byte {
	h := sha256.New()
	for _ , packet :=range list {
		h.Write(packet.Hash)
		h.Write(packet.Signature)
		h.Write(packet.Owner)
	}

	return h.Sum(nil)
}

func packetListHasPacket(packetList []Packet, packetInQuestion Packet) bool {
	for _ , packet := range packetList {
		if equalPackets(packet, packetInQuestion){
			return true
		}
	}
	// fmt.Println("Did not find packet in the list of packets")
	return false
}

func packetListHasPacketHashAndPublicKey(packetList []Packet, packetHash, publicKey []byte) bool {
	for _ , packet := range packetList {
		if string(packet.Hash) == string(packetHash) && string(packet.Owner) == string(publicKey) {
			return true
		}
	}
	// fmt.Println("Did not find a packet in the list of packets that had the designated hash")
	return false
}

func getPacketFromListByHashAndPublicKey(packetList []Packet, packetHash, publicKey []byte) Packet {
	for _ , packet := range packetList {
		if string(packet.Hash) == string(packetHash) && string(packet.Owner) == string(publicKey){
			return packet
		}
	}
	fmt.Println("Did not find a packet in the list of packets that had the designated hash")
	return Packet{}
}

func equalPackets(packet1, packet2 Packet) bool {
	ownerEqual := string(packet1.Owner)     == string(packet2.Owner)
	hashEqual  := string(packet1.Hash)      == string(packet2.Hash)
	sigEqual   := string(packet1.Signature) == string(packet1.Signature)

	return ownerEqual && hashEqual && sigEqual
}












