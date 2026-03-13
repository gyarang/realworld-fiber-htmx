.PHONY: build run dev test test-coverage lint fmt clean e2e

BINARY_NAME=conduit
BUILD_DIR=bin

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) main.go

run:
	go run main.go

dev:
	air

test:
	go test ./... -v

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

lint:
	golangci-lint run

fmt:
	gofmt -w .
	goimports -w .

clean:
	rm -rf $(BUILD_DIR) coverage.out coverage.html tmp

e2e:
	cd e2e && npx playwright test
