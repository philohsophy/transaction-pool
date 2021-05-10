# transaction-pool

![GitHub Workflow Status (branch)](https://img.shields.io/github/workflow/status/philohsophy/transaction-pool/CI/main)

Part of [Dummy-Blockchain project](https://github.com/users/philohsophy/projects/1)

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
