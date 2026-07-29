# ============================================================
# RentOS Backend — Makefile
# ============================================================

BINARY_API    = bin/api
BINARY_WORKER = bin/worker
MODULE        = rentos-backend
GO            = go
GOTEST        = $(GO) test
GOLINT        = golangci-lint

.PHONY: all build build-api build-worker run-api run-worker \
        test test-unit test-cover lint fmt vet \
        migrate migrate-down migrate-status \
        docker-up docker-down docker-logs \
        clean help

# ---- Build ----

all: build

build: build-api build-worker

build-api:
	@echo "→ building api..."
	@mkdir -p bin
	CGO_ENABLED=0 $(GO) build -ldflags="-w -s" -o $(BINARY_API) ./cmd/api

build-worker:
	@echo "→ building worker..."
	@mkdir -p bin
	CGO_ENABLED=0 $(GO) build -ldflags="-w -s" -o $(BINARY_WORKER) ./cmd/worker

# ---- Run ----

run-api:
	$(GO) run ./cmd/api

run-worker:
	$(GO) run ./cmd/worker

# ---- Test ----

test: test-unit

test-unit:
	@echo "→ running unit tests..."
	$(GOTEST) -v -race ./internal/... ./pkg/...

test-cover:
	@echo "→ running tests with coverage..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./internal/... ./pkg/...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "coverage report: coverage.html"

test-pkg:
	@echo "→ running tests for package: $(PKG)"
	$(GOTEST) -v -race $(PKG)

# ---- Code quality ----

lint:
	@echo "→ running linter..."
	$(GOLINT) run ./...

fmt:
	@echo "→ formatting code..."
	$(GO) fmt ./...
	goimports -w .

vet:
	@echo "→ running vet..."
	$(GO) vet ./...

# ---- Database migrations ----

migrate:
	@echo "→ running migrations up..."
	migrate -path ./migrations -database "$$DATABASE_URL" up

migrate-down:
	@echo "→ rolling back last migration..."
	migrate -path ./migrations -database "$$DATABASE_URL" down 1

migrate-status:
	migrate -path ./migrations -database "$$DATABASE_URL" version

migrate-create:
	@read -p "Migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

# ---- Docker ----

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f api worker

docker-build:
	docker compose build --no-cache

# ---- Helpers ----

tidy:
	$(GO) mod tidy

clean:
	rm -rf bin/ coverage.out coverage.html

help:
	@echo ""
	@echo "RentOS Backend — available targets:"
	@echo ""
	@echo "  make build         Build api + worker binaries to bin/"
	@echo "  make run-api       Run API server (dev)"
	@echo "  make run-worker    Run worker process (dev)"
	@echo "  make test          Run unit tests"
	@echo "  make test-cover    Run tests + generate HTML coverage report"
	@echo "  make lint          Run golangci-lint"
	@echo "  make fmt           Format code (gofmt + goimports)"
	@echo "  make migrate       Run DB migrations up"
	@echo "  make migrate-down  Roll back last migration"
	@echo "  make docker-up     Start postgres + redis via docker compose"
	@echo "  make docker-down   Stop docker services"
	@echo "  make clean         Remove binaries and coverage files"
	@echo ""
