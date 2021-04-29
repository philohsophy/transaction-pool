#!/bin/bash
export APP_DB_USERNAME=postgres
export APP_DB_PASSWORD=postgres
export APP_DB_NAME=postgres

go install
~/go/bin/dummy-blockchain-transaction-pool