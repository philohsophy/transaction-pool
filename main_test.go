package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/google/uuid"
	main "github.com/philohsophy/dummy-blockchain-transaction-pool"
)

var a main.App

func TestMain(m *testing.M) {
	a.Initialize(
		os.Getenv("TRANSACTION_POOL_DB_HOST"),
		os.Getenv("TRANSACTION_POOL_DB_PORT"),
		os.Getenv("TRANSACTION_POOL_DB_USERNAME"),
		os.Getenv("TRANSACTION_POOL_DB_PASSWORD"),
		os.Getenv("TRANSACTION_POOL_DB_NAME"))

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

func createTransactions(count int) []main.Transaction {
	if count < 1 {
		count = 1
	}

	recipientAddress := main.Address{Name: "Foo", Street: "FooStreet", HouseNumber: "1", Town: "FooTown"}
	recipientAddressJson, _ := json.Marshal(recipientAddress)
	senderAddress := main.Address{Name: "Bar", Street: "BarStreet", HouseNumber: "1", Town: "BarTown"}
	senderAddressJson, _ := json.Marshal(senderAddress)

	var transactions = make([]main.Transaction, count)
	for i := 0; i < count; i++ {
		var t main.Transaction
		t.Id = uuid.New()
		t.RecipientAddress = recipientAddress
		t.SenderAddress = senderAddress
		t.Value = float32((i + 1.0) * 10)
		transactions[i] = t

		_, err := a.DB.Exec(`
			INSERT INTO transactions
			VALUES($1, $2, $3, $4)
			RETURNING id`,
			t.Id, recipientAddressJson, senderAddressJson, t.Value)

		if err != nil {
			log.Fatal(err)
		}
	}
	return transactions
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

/*
func checkResponseBody(t *testing.T, expectedBody, actualBody interface{}) {}
*/

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/transactions", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != `{"transactions":[]}` {
		t.Errorf("Expected an empty array of transactions. Got %s", body)
	}
}

func TestGetTransactions(t *testing.T) {
	clearTable()
	createTransactions(5)

	t.Run("If no amount is given via queryString params", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/transactions", nil)
		response := executeRequest(req)
		checkResponseCode(t, http.StatusOK, response.Code)

		var m map[string]interface{}
		json.Unmarshal(response.Body.Bytes(), &m)
		// Ref: https://stackoverflow.com/a/21070860
		transactions := m["transactions"].([]interface{})

		expectedAmount := 3
		receivedAmount := len(transactions)

		if receivedAmount != expectedAmount {
			t.Errorf("Expected %v transactions to be returned. Got %v", expectedAmount, receivedAmount)
		}
	})

	t.Run("If amount is given via queryString params", func(t *testing.T) {
		n := 4
		req, _ := http.NewRequest("GET", fmt.Sprintf("/transactions?amount=%d", n), nil)
		response := executeRequest(req)
		checkResponseCode(t, http.StatusOK, response.Code)

		var m map[string]interface{}
		json.Unmarshal(response.Body.Bytes(), &m)
		transactions := m["transactions"].([]interface{})

		expectedAmount := n
		receivedAmount := len(transactions)

		if receivedAmount != expectedAmount {
			t.Errorf("Expected %v transactions to be returned. Got %v", expectedAmount, receivedAmount)
		}
	})

	t.Run("If invalid amount is given via queryString params", func(t *testing.T) {
		invalidAmounts := [5]string{"-1", "0", "", "one", "1xx"}

		for _, invalidAmount := range invalidAmounts {
			req, _ := http.NewRequest("GET", fmt.Sprintf("/transactions?amount=%s", invalidAmount), nil)
			response := executeRequest(req)

			checkResponseCode(t, http.StatusBadRequest, response.Code)
			var m map[string]string
			json.Unmarshal(response.Body.Bytes(), &m)
			expectedErrorMsg := fmt.Sprintf("Invalid amount '%s'", invalidAmount)
			if m["error"] != expectedErrorMsg {
				t.Errorf("Expected the 'error' key of the response to be set to '%s'. Got '%s'", expectedErrorMsg, m["error"])
			}
		}
	})
}

func TestCreateTransaction(t *testing.T) {
	clearTable()

	t.Run("valid transaction", func(t *testing.T) {
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

		var mReq map[string]interface{}
		json.Unmarshal(transactionJson, &mReq)

		var mRes map[string]interface{}
		json.Unmarshal(response.Body.Bytes(), &mRes)

		if !reflect.DeepEqual(mRes["recipientAddress"], mReq["recipientAddress"]) {
			t.Errorf("Expected recipientAddresss to be '%v'. Got '%v'", mReq["recipientAddress"], mRes["recipientAddress"])
		}

		if !reflect.DeepEqual(mRes["senderAddress"], mReq["senderAddress"]) {
			t.Errorf("Expected senderAddress to be '%v'. Got '%v'", mReq["senderAddress"], mRes["senderAddress"])
		}

		if !reflect.DeepEqual(mRes["value"], mReq["value"]) {
			t.Errorf("Expected value to be '%v'. Got '%v'", mReq["value"], mRes["value"])
		}

		id, ok := mRes["id"].(string)
		if !ok {
			t.Errorf("Expected id to be a 'string'. Got '%T'", mRes["id"])
		}
		_, err := uuid.Parse(id)
		if err != nil {
			t.Errorf("Expected id to be an 'UUID'")
		}
	})

	t.Run("invalid transaction", func(t *testing.T) {
		var invalidTransactions = make(map[string][]byte)
		invalidTransactions["recipientAddress"] = []byte(`
			{
				"senderAddress": {
					"name": "Bob",
					"street": "Hauptstrasse",
					"houseNumber": "1",
					"town": "Berlin"
				},
				"value": 100.21
			}`)
		invalidTransactions["senderAddress"] = []byte(`
			{
				"recipientAddress":{
					"name": "Alan",
					"street": "Baker Street",
					"houseNumber": "221B",
					"town": "London"
				},
				"value": 100.21
			}`)
		invalidTransactions["value"] = []byte(`
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
				}
			}`)

		for missingElement, invalidTransaction := range invalidTransactions {
			req, _ := http.NewRequest("POST", "/transactions", bytes.NewBuffer(invalidTransaction))
			req.Header.Set("Content-Type", "application/json; charset=UTF-8")

			response := executeRequest(req)
			checkResponseCode(t, http.StatusBadRequest, response.Code)

			var m map[string]string
			json.Unmarshal(response.Body.Bytes(), &m)
			expectedErrorMsg := fmt.Sprintf("Invalid transaction: missing '%s'", missingElement)
			if m["error"] != expectedErrorMsg {
				t.Errorf("Expected the 'error' key of the response to be set to '%s'. Got '%s'", expectedErrorMsg, m["error"])
			}
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		// missing ',' between recipientAddress and senderAddress --> invalid JSON
		malformedJson := []byte(`
		{
			"recipientAddress":{
				"name": "Alan",
				"street": "Baker Street",
				"houseNumber": "221B",
				"town": "London"
			}
			"senderAddress": {
				"name": "Bob",
				"street": "Hauptstrasse",
				"houseNumber": "1",
				"town": "Berlin"
			},
			"value": 100.21
		}`)
		req, _ := http.NewRequest("POST", "/transactions", bytes.NewBuffer(malformedJson))
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
		response := executeRequest(req)
		checkResponseCode(t, http.StatusBadRequest, response.Code)

		var m map[string]string
		json.Unmarshal(response.Body.Bytes(), &m)
		expectedErrorMsg := "Invalid JSON"
		if m["error"] != expectedErrorMsg {
			t.Errorf("Expected the 'error' key of the response to be set to '%s'. Got '%s'", expectedErrorMsg, m["error"])
		}
	})
}

func TestGetTransaction(t *testing.T) {
	clearTable()

	t.Run("transaction exists (200)", func(t *testing.T) {
		transactions := createTransactions(2)

		for _, transaction := range transactions {
			req, _ := http.NewRequest("GET", "/transactions/"+transaction.Id.String(), nil)
			response := executeRequest(req)
			checkResponseCode(t, http.StatusOK, response.Code)

			transactionJson, _ := json.Marshal(transaction)
			if string(transactionJson) != response.Body.String() {
				t.Errorf("Received transaction does not match!\n\tExpected: '%v'\n\tReceived: '%v'", string(transactionJson), response.Body)
			}
		}
	})

	t.Run("transaction does not exist (404)", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/transactions/"+uuid.New().String(), nil)
		response := executeRequest(req)

		checkResponseCode(t, http.StatusNotFound, response.Code)
		var m map[string]string
		json.Unmarshal(response.Body.Bytes(), &m)
		expectedErrorMsg := "Transaction not found"
		if m["error"] != expectedErrorMsg {
			t.Errorf("Expected the 'error' key of the response to be set to '%s'. Got '%s'", expectedErrorMsg, m["error"])
		}
	})
}

func TestDeleteTransaction(t *testing.T) {
	clearTable()

	t.Run("transaction exists (200)", func(t *testing.T) {
		transactions := createTransactions(1)
		transaction := transactions[0]

		req, _ := http.NewRequest("DELETE", "/transactions/"+transaction.Id.String(), nil)
		response := executeRequest(req)
		checkResponseCode(t, http.StatusOK, response.Code)

		transactionJson, _ := json.Marshal(transaction)
		if string(transactionJson) != response.Body.String() {
			t.Errorf("Received transaction does not match!\n\tExpected: '%v'\n\tReceived: '%v'", string(transactionJson), response.Body)
		}
	})

	t.Run("transaction does not exist (404)", func(t *testing.T) {
		transactionId := uuid.New().String()

		req, _ := http.NewRequest("DELETE", "/transactions/"+transactionId, nil)
		response := executeRequest(req)
		checkResponseCode(t, http.StatusNotFound, response.Code)

		var m map[string]string
		json.Unmarshal(response.Body.Bytes(), &m)
		expectedErrorMsg := "Transaction not found"
		if m["error"] != expectedErrorMsg {
			t.Errorf("Expected the 'error' key of the response to be set to '%s'. Got '%s'", expectedErrorMsg, m["error"])
		}
	})
}
