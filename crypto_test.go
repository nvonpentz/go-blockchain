package main

import (
	"crypto/sha256"
	"testing"
)

func TestKeyGeneration(t *testing.T) {

	keypair := GenerateNewKeypair()

	if len(keypair.Public) > 80 {
		t.Error("Error generating key")
	}
}

func TestKeySigning(t *testing.T) {

	for i := 0; i < 5; i++ {
		keypair := GenerateNewKeypair()

		data := arrayOfBytes(i, 'a')
		hash := SHA256(data)

		signature, err := keypair.Sign(hash)

		if err != nil {

			t.Error("base58 error")

		} else if !SignatureVerify(keypair.Public, signature, hash) {

			t.Error("Signing and verifying error", len(keypair.Public))
		}
	}

}

func SHA256(data []byte) []byte {

	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}
