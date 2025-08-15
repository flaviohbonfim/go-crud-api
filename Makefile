.PHONY: dev build run lint test cover migrate-up migrate-down migrate-create

# Go variables
BINARY_NAME=go-crud-api
CMD_PATH=./cmd/api

# --- Development ---
dev:
	@echo "Running API in dev mode..."
	go run ${CMD_PATH}/main.go

# --- Build ---
build:
	@echo "Building binary..."
	go build -o bin/${BINARY_NAME} ${CMD_PATH}/main.go

run:
	@echo "Running binary..."
	./bin/${BINARY_NAME}

# --- Quality ---
lint:
	@echo "Running linter..."
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run ./...

test:
	@echo "Running tests..."
	go test ./... -v

cover:
	@echo "Running tests with coverage..."
	go test ./... -coverprofile=coverage.out

# --- Database ---
migrate-up:
	@echo "Applying migrations..."
	# migrate -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -path migrations up

migrate-down:
	@echo "Reverting migrations..."
	# migrate -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -path migrations down

migrate-create:
	@read -p "Enter migration name: " name; \
	# migrate create -ext sql -dir migrations -seq $$name
