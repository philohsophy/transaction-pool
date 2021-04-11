package models

type Address struct {
	Name        string `json:"name"`
	Street      string `json:"street"`
	HouseNumber int    `json:"houseNumber"`
	Town        string `json:"town"`
}

type Transaction struct {
	Id              int     `json:"id"`
	RecipentAddress Address `json:"recipientAddress"`
	SenderAddress   Address `json:"senderAddress"`
	Value           string  `json:"value"`
}
