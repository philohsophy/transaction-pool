package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	models "github.com/philohsophy/dummy-blockchain-transaction-pool/transaction"
)

var transactions []models.Transaction

func initTransactionPool() {
	aR := models.Address{Name: "Alan", Street: "Baker Street", HouseNumber: 1, Town: "London"}
	aS := models.Address{Name: "Bob", Street: "Hauptstrasse", HouseNumber: 11, Town: "Berlin"}

	transactions = append(transactions, models.Transaction{Id: uuid.New(), RecipentAddress: aR, SenderAddress: aS, Value: "123"})
	transactions = append(transactions, models.Transaction{Id: uuid.New(), RecipentAddress: aR, SenderAddress: aS, Value: "456"})
	transactions = append(transactions, models.Transaction{Id: uuid.New(), RecipentAddress: aR, SenderAddress: aS, Value: "789"})
}

func ListTransactions(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(transactions)
}

func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	transaction := models.Transaction{}
	_ = json.NewDecoder(r.Body).Decode(&transaction)
	transaction.Id = uuid.New()
	transactions = append(transactions, transaction)

	json.NewEncoder(w).Encode(transaction)
}

func GetTransaction(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	for _, transaction := range transactions {
		if transaction.Id.String() == params["id"] {
			json.NewEncoder(w).Encode(transaction)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func main() {
	initTransactionPool()

	r := mux.NewRouter()
	r.HandleFunc("/transactions", ListTransactions).Methods("GET")
	r.HandleFunc("/transactions", CreateTransaction).Methods("POST")
	r.HandleFunc("/transactions/{id}", GetTransaction).Methods("GET")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}
