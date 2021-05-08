package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	models "github.com/philohsophy/dummy-blockchain-models"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(host, port, user, password, dbname string) {
	portInt, _ := strconv.Atoi(port)
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, portInt, user, password, dbname)
	a.connectToDatabase(connectionString)
	a.initializeDatabase()

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) connectToDatabase(connectionString string) {
	fmt.Print("Connecting to database...")
	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(" ok")
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS transactions
(
	id UUID PRIMARY KEY,
	recipient_address JSONB NOT NULL,
	sender_address JSONB NOT NULL,
	value NUMERIC(10,2) NOT NULL DEFAULT 0.00
)`

func (a *App) initializeDatabase() {
	fmt.Print("Initializing database...")
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
	fmt.Println(" ok")
}

func (a *App) initializeRoutes() {
	fmt.Print("Initializing routes...")
	a.Router.HandleFunc("/transactions", a.getTransactions).Methods("GET")
	a.Router.HandleFunc("/transactions", a.createTransaction).Methods("POST")
	a.Router.HandleFunc("/transactions/{id:[a-z0-9-]+}", a.getTransaction).Methods("GET")
	a.Router.HandleFunc("/transactions/{id:[a-z0-9-]+}", a.deleteTransaction).Methods("DELETE")
	fmt.Println(" ok")
}

func (a *App) Run(addr string) {
	log.Printf("Serving on %s", addr)
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJson(w, code, map[string]string{"error": message})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) getTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["id"])
	if err != nil {
		fmt.Println("Error: ", err.Error())
		respondWithError(w, http.StatusBadRequest, "Invalid transaction ID")
		return
	}

	t := &Transaction{&models.Transaction{Id: id}} // see https://stackoverflow.com/a/60518886
	if err := t.getTransaction(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Transaction not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJson(w, http.StatusOK, t)
}

func (a *App) createTransaction(w http.ResponseWriter, r *http.Request) {
	var t Transaction
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&t); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	defer r.Body.Close()

	if err := t.createTransaction(a.DB); err != nil {
		switch err.(type) {
		case *InvalidTransactionError:
			respondWithError(w, http.StatusBadRequest, err.Error())
		case *pq.Error:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJson(w, http.StatusCreated, t)
}

func (a *App) deleteTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["id"])
	if err != nil {
		fmt.Println("Error: ", err.Error())
		respondWithError(w, http.StatusBadRequest, "Invalid transaction ID")
		return
	}

	t := &Transaction{&models.Transaction{Id: id}} // see https://stackoverflow.com/a/60518886
	if err := t.deleteTransaction(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Transaction not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJson(w, http.StatusOK, t)
}

func (a *App) getTransactions(w http.ResponseWriter, r *http.Request) {
	var amount int = 3 // Default

	query := r.URL.Query()
	amountQueryString, present := query["amount"]
	if present && len(amountQueryString) > 0 {
		amountRequested, err := strconv.Atoi(amountQueryString[0])
		if err != nil || amountRequested <= 0 {
			respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid amount '%s'", amountQueryString[0]))
			return
		}
		amount = amountRequested
	}

	transactions, err := getTransactions(a.DB, amount)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	payload := struct {
		Transactions []models.Transaction `json:"transactions"`
	}{
		Transactions: transactions,
	}
	respondWithJson(w, http.StatusOK, payload)
}
