package main

import (
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"
	models "github.com/philohsophy/dummy-blockchain-models"
)

type InvalidTransactionError struct {
	err string
}

func (e *InvalidTransactionError) Error() string {
	return e.err
}

type Transaction struct {
	*models.Transaction
}

func (t *Transaction) createTransaction(db *sql.DB) error {
	t.Id = uuid.New()
	if !t.IsValid() {
		return &InvalidTransactionError{err: "Invalid transaction"}
	}

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
	return db.QueryRow(`
		SELECT * FROM transactions
		WHERE id=$1`,
		t.Id).Scan(&t.Id, &t.RecipientAddress, &t.SenderAddress, &t.Value)
}

func (t *Transaction) deleteTransaction(db *sql.DB) error {
	// Ref: https://www.calhoun.io/updating-and-deleting-postgresql-records-using-gos-sql-package/
	return db.QueryRow(`
		DELETE FROM transactions
		WHERE id=$1
		RETURNING id, recipient_address, sender_address, value`,
		t.Id).Scan(&t.Id, &t.RecipientAddress, &t.SenderAddress, &t.Value)
}

func getTransactions(db *sql.DB, count int) ([]models.Transaction, error) {
	// TODO: check if count is nil
	rows, err := db.Query(`
		SELECT * FROM transactions
		LIMIT $1`,
		count)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	transactions := []models.Transaction{}

	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.Id, &t.RecipientAddress, &t.SenderAddress, &t.Value); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}
