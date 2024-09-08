include .env

SRC=src/main/cmd/main.go

.PHONY: all
all:
	make start & PID=$$; \
	sleep 1; \
	make migrate-up; \
	trap 'make drop-table; kill $$PID' EXIT; \
	wait

.PHONY: setup
setup:
	@echo "Installing dependencies..."
	@go mod download
	@go mod download github.com/gorilla/mux
	@go mod download github.com/lib/pq
	@echo "Dependencies installed"

.PHONY: start
start:
	@go run $(SRC)

.PHONY: migrate-up
migrate-up:
	@migrate -path migrations -database "$(POSTGRES_CONN)" up

.PHONY: migrate-down
migrate-down:
	@migrate -path migrations -database "$(POSTGRES_CONN)" down

.PHONY: drop-table
drop-db:
	@migrate -path migrations -database "$(POSTGRES_CONN)" drop -f