package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	models "github.com/philohsophy/dummy-blockchain-transaction-pool/transaction"
)

var transactionQueue []models.Transaction

func initTransactionPool() {
	aR := models.Address{Name: "Alan", Street: "Baker Street", HouseNumber: 1, Town: "London"}
	aS := models.Address{Name: "Bob", Street: "Hauptstrasse", HouseNumber: 11, Town: "Berlin"}

	transactionQueue = append(transactionQueue, models.Transaction{Id: 1, RecipentAddress: aR, SenderAddress: aS, Value: "123"})
	transactionQueue = append(transactionQueue, models.Transaction{Id: 2, RecipentAddress: aR, SenderAddress: aS, Value: "456"})
	transactionQueue = append(transactionQueue, models.Transaction{Id: 3, RecipentAddress: aR, SenderAddress: aS, Value: "789"})
}

func ListTransactions(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(transactionQueue)
}

func GetTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid transaction Id"))
	}

	for _, transaction := range transactionQueue {
		if transaction.Id == id {
			json.NewEncoder(w).Encode(transaction)
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func main() {
	initTransactionPool()

	r := mux.NewRouter()
	r.HandleFunc("/transactions", ListTransactions)
	r.HandleFunc("/transactions/{id}", GetTransaction)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}
