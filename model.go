package main

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
)

type Address struct {
	Name        string `json:"name"`
	Street      string `json:"street"`
	HouseNumber string `json:"houseNumber"`
	Town        string `json:"town"`
}

type Transaction struct {
	Id               uuid.UUID `json:"id"`
	RecipientAddress Address   `json:"recipientAddress"`
	SenderAddress    Address   `json:"senderAddress"`
	Value            float32   `json:"value"`
}

func (t *Transaction) createTransaction(db *sql.DB) error {
	t.Id = uuid.New()
	recipientAddressJson, _ := json.Marshal(t.RecipientAddress)
	senderAddressJson, _ := json.Marshal(t.SenderAddress)

	_, err := db.Exec(`
		INSERT INTO transactions
		VALUES($1, $2, $3, $4)
		RETURNING id`,
		t.Id, recipientAddressJson, senderAddressJson, t.Value)

	if err != nil {
		return err
	}

	return nil
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
