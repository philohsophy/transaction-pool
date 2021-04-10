package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", TransactionHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func TransactionHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello there")
}
