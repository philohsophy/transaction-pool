# transaction-pool

![GitHub Workflow Status (branch)](https://img.shields.io/github/workflow/status/philohsophy/transaction-pool/CI/main)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/philohsophy/transaction-pool)

Part of [Blockchain-Demo project](https://github.com/philohsophy/blockchain-demo)

## Outline

Backend Service written in GO which:

- provides a HTTP/REST-Interface for managing transactions
- uses a PostgreSQL database for persisting the transactions

## How to run

```bash
# 1. start postgres db
./scripts/run_database.sh

# 2. run app
./scripts/run.sh

# 3. i.e. create new transaction
./scripts/create_transaction.sh
```

Or import & use Postman Collection located in ```./doc/postman/```
