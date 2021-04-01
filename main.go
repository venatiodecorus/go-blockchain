package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var bc *Blockchain

func main() {
	var err error
	bc, err = NewBlockchain("newChain")
	if err != nil {
		panic("error creating chain")
	}

	r := mux.NewRouter()
	r.HandleFunc("/lastBlock", LastBlockHandler)

	r.HandleFunc("/wallets", ListWalletsHandler).Methods("GET")
	r.HandleFunc("/wallets", CreateWalletHandler).Methods("POST")
	r.HandleFunc("/wallets/{address}", WalletBalanceHandler).Methods("GET")

	r.HandleFunc("/transaction", TransactionHandler).Methods("POST")

	http.Handle("/", r)

	http.ListenAndServe(":8080", r)
}

func TransactionHandler(w http.ResponseWriter, r *http.Request) {
	amount, err := strconv.Atoi(r.FormValue("amount"))
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		bc.NewTrx(r.FormValue("to"), r.FormValue("from"), amount)
	}
}

func LastBlockHandler(w http.ResponseWriter, r *http.Request) {
	lastBlock, err := bc.LastBlock()
	if err != nil {
		w.Write([]byte("ERROR"))
	} else {
		json.NewEncoder(w).Encode(lastBlock)
	}
}

func ListWalletsHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(bc.Wallets)
}

func CreateWalletHandler(w http.ResponseWriter, r *http.Request) {
	address, err := bc.CreateWallet()
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		json.NewEncoder(w).Encode(address)
	}
}

func WalletBalanceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	balance, err := bc.GetWalletBalance(address)
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		json.NewEncoder(w).Encode(balance)
	}
}
