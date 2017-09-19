package main

/*
The majority of this code was written by Jorge Izquierdo for his implementation of a blockchain written in Go
https://github.com/izqui/blockchain

His code offered an elegant abstraction big number types used for ECDSA from the Go crypto library. 
*/

import (
	"fmt"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"math/big"
	"github.com/tv42/base58"
)

// Key generation with proof of work
type Keypair struct {
	Public  []byte `json:"public"`  // base58 (x y)
	Private []byte `json:"private"` // d (base58 encoded)
}

func GenerateNewKeypair() *Keypair {

	pk, _ := ecdsa.GenerateKey(elliptic.P224(), rand.Reader)

	b := bigJoin(28, pk.PublicKey.X, pk.PublicKey.Y)

	public := base58.EncodeBig([]byte{}, b)
	private := base58.EncodeBig([]byte{}, pk.D)

	kp := Keypair{Public: public, Private: private}

	return &kp
}

func (k *Keypair) Sign(hash []byte) ([]byte, error) {

	d, err := base58.DecodeToBig(k.Private)
	if err != nil {
		return nil, err
	}

	b, _ := base58.DecodeToBig(k.Public)

	pub := splitBig(b, 2)
	x, y := pub[0], pub[1]

	key := ecdsa.PrivateKey{ecdsa.PublicKey{elliptic.P224(), x, y}, d}

	r, s, _ := ecdsa.Sign(rand.Reader, &key, hash)

	return base58.EncodeBig([]byte{}, bigJoin(28, r, s)), nil
}

func SignatureVerify(publicKey, sig, hash []byte) bool {

	b, _ := base58.DecodeToBig(publicKey)
	publ := splitBig(b, 2)
	x, y := publ[0], publ[1]

	b, _ = base58.DecodeToBig(sig)
	sigg := splitBig(b, 2)
	r, s := sigg[0], sigg[1]

	pub := ecdsa.PublicKey{elliptic.P224(), x, y}

	return ecdsa.Verify(&pub, hash, r, s)
}

func PrintKeys(keys Keypair){
	fmt.Printf("Public: %v\n", string(keys.Public))
	fmt.Printf("Private: %v\n", string(keys.Private))
}

func bigJoin(expectedLen int, bigs ...*big.Int) *big.Int {

	bs := []byte{}
	for i, b := range bigs {

		by := b.Bytes()
		dif := expectedLen - len(by)
		if dif > 0 && i != 0 {

			by = append(arrayOfBytes(dif, 0), by...)
		}

		bs = append(bs, by...)
	}

	b := new(big.Int).SetBytes(bs)

	return b
}

func splitBig(b *big.Int, parts int) []*big.Int {

	bs := b.Bytes()
	if len(bs)%2 != 0 {
		bs = append([]byte{0}, bs...)
	}

	l := len(bs) / parts
	as := make([]*big.Int, parts)

	for i, _ := range as {

		as[i] = new(big.Int).SetBytes(bs[i*l : (i+1)*l])
	}

	return as

}

func arrayOfBytes(i int, b byte) (p []byte) {

	for i != 0 {

		p = append(p, b)
		i--
	}
	return
}