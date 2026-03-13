.PHONY: build run dev test test-coverage lint fmt clean e2e check coverage-check

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

check: lint test

COVERAGE_PKGS=./cmd/web/controller ./cmd/web/controller/htmx ./cmd/web/model ./internal

coverage-check:
	@echo "Checking test coverage per package (minimum 80%)..."
	@FAIL=0; \
	for pkg in $(COVERAGE_PKGS); do \
		RESULT=$$(go test -cover $$pkg 2>&1); \
		COV=$$(echo "$$RESULT" | grep -oE '[0-9]+\.[0-9]+%' | tr -d '%'); \
		if [ -z "$$COV" ]; then \
			echo "  SKIP: $$pkg (no test files)"; \
		elif [ $$(echo "$$COV < 80.0" | bc) -eq 1 ]; then \
			echo "  FAIL: $$pkg — $${COV}%"; \
			FAIL=1; \
		else \
			echo "  OK:   $$pkg — $${COV}%"; \
		fi; \
	done; \
	if [ $$FAIL -eq 1 ]; then \
		echo "FAIL: One or more packages below 80% threshold"; \
		exit 1; \
	else \
		echo "OK: All packages meet 80% threshold"; \
	fi
