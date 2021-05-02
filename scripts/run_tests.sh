#!/bin/bash
export TRANSACTION_POOL_DB_HOST=localhost
export TRANSACTION_POOL_DB_PORT=5432
export TRANSACTION_POOL_DB_USERNAME=postgres
export TRANSACTION_POOL_DB_PASSWORD=postgres
export TRANSACTION_POOL_DB_NAME=postgres

# i.e. ./scripts/run_test.sh -run Create
go test -v $1 $2