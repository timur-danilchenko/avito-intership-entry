include .env

SRC=cmd/main.go
MIGRATIONS_PATH=migrations

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
	@go get -u github.com/golang-migrate/migrate/v4
	@echo "Migration utility installed"
	

.PHONY: start
start:
	@go run $(SRC)

.PHONY: migrate-up
migrate-up:
	@migrate -path ${MIGRATIONS_PATH} -database "$(POSTGRES_CONN)" up

.PHONY: migrate-down
migrate-down:
	@migrate -path ${MIGRATIONS_PATH} -database "$(POSTGRES_CONN)" down 1

.PHONY: drop-table
drop-db:
	@migrate -path ${MIGRATIONS_PATH} -database "$(POSTGRES_CONN)" drop -f

.PHONY: dbshell
dbshell:
	@psql -U "${POSTGRES_USER}" -d "${POSTGRES_DATABASE}"