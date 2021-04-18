export APP_DB_USERNAME=postgres
export APP_DB_PASSWORD=postgres
export APP_DB_NAME=postgres

docker run -d --rm \
-p 5432:5432 \
-e POSTGRES_PASSWORD=postgres \
--name postgres \
postgres:13.2-alpine