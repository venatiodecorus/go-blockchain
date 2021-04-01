package main

import (
	"crypto/sha256"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

var bc *Blockchain

func main() {
	// fmt.Println("lol")

	var err error
	bc, err = NewBlockchain("newChain")
	if err != nil {
		panic("error creating chain")
	}

	r := mux.NewRouter()
	r.HandleFunc("/", TransactionHandler)
	r.HandleFunc("/lastBlock", LastBlockHandler)
	http.Handle("/", r)

	h := sha256.New()
	h.Write([]byte("tony"))
	to := h.Sum(nil)

	h = sha256.New()
	h.Write([]byte("francis"))
	from := h.Sum(nil)
	bc.NewTrx(to, from, 20)
	bc.NewTrx(to, from, 30)
	bc.NewTrx(from, to, 10)

	proof := sha256.New()
	proof.Write([]byte("proof string"))

	bc.NewBlock(proof)

	http.ListenAndServe(":8080", r)

	// mallory := sha256.New()
	// mallory.Write([]byte("mallory"))

	// frank := sha256.New()
	// frank.Write([]byte("frank"))

	// bc.NewTrx(mallory, frank, 100)
	// bc.NewTrx(frank, mallory, 50)
	// bc.NewBlock(proof)
}

func TransactionHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Println(r.GetBody())
	w.Write([]byte("hello!"))
}

func LastBlockHandler(w http.ResponseWriter, r *http.Request) {
	lastBlock, err := bc.LastBlock()
	if err != nil {
		w.Write([]byte("ERROR"))
	} else {
		json.NewEncoder(w).Encode(lastBlock)
	}
}
