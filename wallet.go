package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
)

type Wallet struct {
	Address string
	key     rsa.PrivateKey
}

func NewWallet() (*Wallet, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return &Wallet{}, err
	}

	// Generate address from public key
	publicBytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		return &Wallet{}, err
	}

	h := sha256.New()
	h.Write(publicBytes)
	address := hex.EncodeToString(h.Sum(nil))

	return &Wallet{string(address), *key}, nil
}
