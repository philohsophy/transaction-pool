package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	models "github.com/philohsophy/dummy-blockchain-transaction-pool/transaction"
)

type myHandler struct{}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "URL: "+r.URL.String())
}

func Tmp(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "version 2")
}

func TransactionHandler(w http.ResponseWriter, r *http.Request) {
	aR := models.Address{Name: "Alan", Street: "Baker Street", HouseNumber: 1, Town: "London"}
	aS := models.Address{Name: "Bob", Street: "Hauptstrasse", HouseNumber: 11, Town: "Berlin"}

	transactions := []models.Transaction{
		{RecipentAddress: aR, SenderAddress: aS, Value: "123"},
		{RecipentAddress: aR, SenderAddress: aS, Value: "456"},
		{RecipentAddress: aR, SenderAddress: aS, Value: "789"},
	}

	json.NewEncoder(w).Encode(transactions)
}

func main() {
	mux := http.NewServeMux()

	mux.Handle("/", &myHandler{})

	mux.HandleFunc("/tmp", Tmp)

	mux.HandleFunc("/transactions", TransactionHandler)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
