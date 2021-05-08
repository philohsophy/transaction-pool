#!/bin/bash
export TRANSACTION_POOL_DB_HOST=localhost
export TRANSACTION_POOL_DB_PORT=5432
export TRANSACTION_POOL_DB_USERNAME=postgres
export TRANSACTION_POOL_DB_PASSWORD=postgres
export TRANSACTION_POOL_DB_NAME=postgres

go install
~/go/bin/transaction-pool