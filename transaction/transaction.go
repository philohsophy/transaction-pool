package models

import "github.com/google/uuid"

type Address struct {
	Name        string `json:"name"`
	Street      string `json:"street"`
	HouseNumber int    `json:"houseNumber"`
	Town        string `json:"town"`
}

type Transaction struct {
	Id               uuid.UUID `json:"id"`
	RecipientAddress Address   `json:"recipientAddress"`
	SenderAddress    Address   `json:"senderAddress"`
	Value            string    `json:"value"`
}
