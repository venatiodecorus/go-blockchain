package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"time"
)

type Blockchain struct {
	Name    string
	chain   []Block
	pending []Trx
}

type Block struct {
	Timestamp    time.Time `json:"timestamp"`
	Transactions []Trx     `json:"transactions"`
	Proof        hash.Hash `json:"proof"`
	Previous     hash.Hash `json:"previous"`
}

type Trx struct {
	To     []byte
	From   []byte
	Amount int
}

func NewBlockchain(name string) (*Blockchain, error) {
	bc := &Blockchain{name, make([]Block, 0), make([]Trx, 0)}

	// Genesis block
	proof := sha256.New()
	proof.Write([]byte("proof"))
	bc.NewBlock(proof)

	return bc, nil
}

func (bc *Blockchain) GetName() string {
	return bc.Name
}

func (bc *Blockchain) NewBlock(proof hash.Hash) {
	lastBlock, err := bc.LastBlock()
	var previous hash.Hash
	if err != nil {
		h := sha256.New()
		h.Write([]byte("genesis block"))
		previous = h
	} else {
		previous = bc.Hash(lastBlock)
	}

	block := &Block{time.Now(), bc.pending, proof, previous}
	bc.pending = nil
	bc.chain = append(bc.chain, *block)
}

func (bc *Blockchain) LastBlock() (Block, error) {
	if len(bc.chain) == 0 {
		return Block{}, errors.New("no last block")
		// return Block{time.Now(), nil, nil, nil}
	}
	return bc.chain[len(bc.chain)-1], nil
}

func (bc *Blockchain) NewTrx(to []byte, from []byte, amount int) {
	trx := &Trx{to, from, amount}
	bc.pending = append(bc.pending, *trx)
}

func (bc *Blockchain) Hash(block Block) hash.Hash {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", block)))
	return h
}
