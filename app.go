package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(user, password, dbname string) {
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/transactions", a.getTransactions).Methods("GET")
	a.Router.HandleFunc("/transactions", a.createTransaction).Methods("POST")
	a.Router.HandleFunc("/transactions/{id:[a-z0-9-]+}", a.getTransaction).Methods("GET")
	a.Router.HandleFunc("/transactions/{id:[a-z0-9-]+}", a.deleteTransaction).Methods("DELETE")
}

func (a *App) Run(addr string) {
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

	t := Transaction{Id: id}
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
		switch err {
		case err.(*InvalidTransactionError):
			respondWithError(w, http.StatusBadRequest, err.Error())
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

	t := Transaction{Id: id}
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
	count := 1

	transactions, err := getTransactions(a.DB, count)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	payload := struct {
		Transactions []Transaction `json:"transactions"`
	}{
		Transactions: transactions,
	}
	respondWithJson(w, http.StatusOK, payload)
}
