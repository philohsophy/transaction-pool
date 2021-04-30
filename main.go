package main

import "os"

func main() {
	a := App{}
	a.Initialize(
		os.Getenv("TRANSACTION_POOL_DB_HOST"),
		os.Getenv("TRANSACTION_POOL_DB_PORT"),
		os.Getenv("TRANSACTION_POOL_DB_USERNAME"),
		os.Getenv("TRANSACTION_POOL_DB_PASSWORD"),
		os.Getenv("TRANSACTION_POOL_DB_NAME"))

	a.Run(":8010")
}
