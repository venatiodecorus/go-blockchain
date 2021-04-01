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
	Wallets []Wallet
}

type Block struct {
	Timestamp    time.Time `json:"timestamp"`
	Transactions []Trx     `json:"transactions"`
	Proof        hash.Hash `json:"proof"`
	Previous     hash.Hash `json:"previous"`
}

type Trx struct {
	To     string
	From   string
	Amount int
}

func NewBlockchain(name string) (*Blockchain, error) {
	bc := &Blockchain{name, make([]Block, 0), make([]Trx, 0), make([]Wallet, 0)}

	// Genesis block
	proof := sha256.New()
	proof.Write([]byte("proof"))
	bc.NewBlock(proof)

	// Mint wallet
	mint, err := bc.CreateWallet()
	if err != nil {
		return &Blockchain{}, err
	}

	bc.NewTrx(mint, "", 1000000000000)

	ticker := time.NewTicker(500 * time.Millisecond)
	go func() {
		for {
			select {
			case <-ticker.C:
				proof := sha256.New()
				proof.Write([]byte("proof string"))
				bc.NewBlock(proof)
			}
		}
	}()

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

func (bc *Blockchain) NewTrx(to string, from string, amount int) {
	trx := &Trx{to, from, amount}
	bc.pending = append(bc.pending, *trx)
}

func (bc *Blockchain) Hash(block Block) hash.Hash {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", block)))
	return h
}

func (bc *Blockchain) CreateWallet() (string, error) {
	w, err := NewWallet()
	if err != nil {
		return "", err
	}

	bc.Wallets = append(bc.Wallets, *w)
	return w.Address, nil
}

func (bc *Blockchain) GetWalletBalance(address string) (int, error) {
	balance := 0
	for _, block := range bc.chain {
		for _, trx := range block.Transactions {
			if trx.To == address {
				balance += trx.Amount
			}
			if trx.From == address {
				balance -= trx.Amount
			}
		}
	}

	return balance, nil
}
