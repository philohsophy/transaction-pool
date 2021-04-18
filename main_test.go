package main_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	main "github.com/philohsophy/dummy-blockchain-transaction-pool"
)

var a main.App

func TestMain(m *testing.M) {
	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"))

	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM transactions")
	a.DB.Exec("ALTER SEQUENCE transactions_id_seq RESTART WITH 1")
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS transactions
(
	id UUID PRIMARY KEY,
	recipient_address JSONB NOT NULL,
	sender_address JSONB NOT NULL,
	value NUMERIC(10,2) NOT NULL DEFAULT 0.00
)`

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expectedCode, actualCode int) {
	if expectedCode != actualCode {
		t.Errorf("Expected response code %d. Got %d\n", expectedCode, actualCode)
	}
}

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/transactions", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array of transactions. Got %s", body)
	}
}

func TestGetNonExistentTransaction(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/transactions/1337", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)
	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	expectedErrorMsg := "Transaction not found"
	if m["error"] != expectedErrorMsg {
		t.Errorf("Expected the 'error' key of the response to be set to '%s'. Got '%s'", expectedErrorMsg, m["error"])
	}
}

func TestCreateTransaction(t *testing.T) {
	clearTable()

	transactionJson := []byte(`
		{
			"recipientAddress":{
				"name": "Alan",
				"street": "Baker Street",
				"houseNumber": "221B",
				"town": "London"
			},
			"senderAddress": {
				"name": "Bob",
				"street": "Hauptstrasse",
				"houseNumber": "1",
				"town": "Berlin"
			},
			"value": 100.21
		}`)
	req, _ := http.NewRequest("POST", "/transactions", bytes.NewBuffer(transactionJson))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	log.Println(m)

	// TODO: Define Assertions
}
