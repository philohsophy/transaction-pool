#/bin/bash
HOST=$(hostname -I)

docker run -it --rm \
-p 8010:8010 \
-e TRANSACTION_POOL_DB_HOST=$HOST \
-e TRANSACTION_POOL_DB_PORT=5432 \
-e TRANSACTION_POOL_DB_USERNAME=postgres \
-e TRANSACTION_POOL_DB_PASSWORD=postgres \
-e TRANSACTION_POOL_DB_NAME=postgres \
--name transaction-pool \
transaction-pool:0.0.1