package main

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type Address struct {
	Name        string `json:"name"`
	Street      string `json:"street"`
	HouseNumber int    `json:"houseNumber"`
	Town        string `json:"town"`
}

type Transaction struct {
	Id              uuid.UUID `json:"id"`
	RecipentAddress Address   `json:"recipientAddress"`
	SenderAddress   Address   `json:"senderAddress"`
	Value           float32   `json:"value"`
}

func (t *Transaction) createTransaction(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (t *Transaction) getTransaction(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (t *Transaction) deleteTransaction(db *sql.DB) error {
	return errors.New("Not implemented")
}

func getTransactions(db *sql.DB, start, count int) ([]Transaction, error) {
	return nil, errors.New("Not implemented")
}
