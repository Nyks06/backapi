# Fill all required env vars

export POSTGRESQL_DATASOURCE="dbname=backapi_test user=postgres sslmode=disable host=localhost"
createdb -U postgres -h localhost -p 5432 backapi_test

goose up
go test ./...

dropdb -U postgres -h localhost -p 5432 backapi_test || true