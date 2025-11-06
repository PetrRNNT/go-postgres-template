.PHONY: setup run migrate-up migrate-down generate test db-up db-down

setup:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

run:
	go run cmd/app/main.go

migrate-up:
	migrate -path ./migrations -database "postgresql://user:password@localhost:5432/mydb?sslmode=disable" up

migrate-down:
	migrate -path ./migrations -database "postgresql://user:password@localhost:5432/mydb?sslmode=disable" down

generate:
	sqlc generate

test:
	go test -v ./...

db-up:
	docker-compose up -d

db-down:
	docker-compose down -v
