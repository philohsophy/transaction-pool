package main_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/uuid"
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

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS transactions
(
	id UUID PRIMARY KEY,
	recipient_address JSONB NOT NULL,
	sender_address JSONB NOT NULL,
	value NUMERIC(10,2) NOT NULL DEFAULT 0.00
)`

func clearTable() {
	a.DB.Exec("DELETE FROM transactions *")
}

func createTransactions(count int) []uuid.UUID {
	if count < 1 {
		count = 1
	}

	recipientAddress := main.Address{Name: "Foo", Street: "FooStreet", HouseNumber: "1", Town: "FooTown"}
	recipientAddressJson, _ := json.Marshal(recipientAddress)
	senderAddress := main.Address{Name: "Bar", Street: "BarStreet", HouseNumber: "1", Town: "BarTown"}
	senderAddressJson, _ := json.Marshal(senderAddress)

	var transactionIds = make([]uuid.UUID, count)
	for i := 0; i < count; i++ {
		transactionIds[i] = uuid.New()
		_, err := a.DB.Exec(`
			INSERT INTO transactions
			VALUES($1, $2, $3, $4)
			RETURNING id`,
			transactionIds[i], recipientAddressJson, senderAddressJson, (i+1.0)*10)

		if err != nil {
			log.Fatal(err)
		}
	}
	return transactionIds
}

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

	if body := response.Body.String(); body != `{"transactions":[]}` {
		t.Errorf("Expected an empty array of transactions. Got %s", body)
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

func TestGetNonExistentTransaction(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/transactions/"+uuid.New().String(), nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)
	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	expectedErrorMsg := "Transaction not found"
	if m["error"] != expectedErrorMsg {
		t.Errorf("Expected the 'error' key of the response to be set to '%s'. Got '%s'", expectedErrorMsg, m["error"])
	}
}

func TestGetTransaction(t *testing.T) {
	clearTable()
	transactionIds := createTransactions(2)

	for _, transactionId := range transactionIds {
		req, _ := http.NewRequest("GET", "/transactions/"+transactionId.String(), nil)
		response := executeRequest(req)
		checkResponseCode(t, http.StatusOK, response.Code)
	}
}
