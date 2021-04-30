package main

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"reflect"

	"github.com/google/uuid"
)

type Address struct {
	Name        string `json:"name"`
	Street      string `json:"street"`
	HouseNumber string `json:"houseNumber"`
	Town        string `json:"town"`
}

// Make the Address struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
// Ref: https://www.alexedwards.net/blog/using-postgresql-jsonb
func (a Address) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Make the Address struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
// Ref: https://www.alexedwards.net/blog/using-postgresql-jsonb
func (a *Address) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

type Transaction struct {
	Id               uuid.UUID `json:"id"`
	RecipientAddress Address   `json:"recipientAddress"`
	SenderAddress    Address   `json:"senderAddress"`
	Value            float32   `json:"value"`
}

type InvalidTransactionError struct {
	err string
}

func (e *InvalidTransactionError) Error() string {
	return e.err
}

func (t *Transaction) createTransaction(db *sql.DB) error {
	t.Id = uuid.New()
	recipientAddressJson, _ := json.Marshal(t.RecipientAddress)
	senderAddressJson, _ := json.Marshal(t.SenderAddress)

	if reflect.ValueOf(t.RecipientAddress).IsZero() {
		return &InvalidTransactionError{err: "Invalid transaction: missing 'recipientAddress'"}
	}
	if reflect.ValueOf(t.SenderAddress).IsZero() {
		return &InvalidTransactionError{err: "Invalid transaction: missing 'senderAddress'"}
	}
	if reflect.ValueOf(t.Value).IsZero() {
		return &InvalidTransactionError{err: "Invalid transaction: missing 'value'"}
	}

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

func getTransactions(db *sql.DB, count int) ([]Transaction, error) {
	// TODO: check if count is nil
	rows, err := db.Query(`
		SELECT * FROM transactions
		LIMIT $1`,
		count)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	transactions := []Transaction{}

	for rows.Next() {
		var t Transaction
		if err := rows.Scan(&t.Id, &t.RecipientAddress, &t.SenderAddress, &t.Value); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}
