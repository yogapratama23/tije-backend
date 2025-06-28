include .env

start-api:
	@go run ./cmd/api/

start-receiver:
	@go run ./cmd/receiver/

create-migration:
	@migrate create -ext sql -dir database/migrations $(name)
migrate-up:
	@migrate -database "${DATABASE_URL}" -path database/migrations up
migrate-down:
	@migrate -database "postgresql://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_DATABASE}?sslmode=disable" -path database/migrations down
