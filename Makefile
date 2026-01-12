# CONFIG
APP_NAME=ecommerce
CMD_DIR=cmd/api
BIN_DIR=bin
MAIN_FILE=$(CMD_DIR)/main.go

GO=go
GOFLAGS=-v
ENV?=dev

# DEFAULT
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make run          Run API locally"
	@echo "  make build        Build binary"
	@echo "  make test         Run tests"
	@echo "  make test-cover   Run tests with coverage"
	@echo "  make clean        Clean build artifacts"

# RUN
.PHONY: run
run:
	ENV=$(ENV) $(GO) run $(MAIN_FILE)

# BUILD
.PHONY: build
build:
	mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	$(GO) build -o $(BIN_DIR)/$(APP_NAME) $(MAIN_FILE)

# TESTING
PHONY: test
test:
	$(GO) test ./...

.PHONY: test-cover
test-cover:
	$(GO) test ./... -coverprofile=coverage.out
	$(GO) tool cover -html=coverage.out

# CLEAN
.PHONY: clean
clean:
	rm -rf $(BIN_DIR)
	rm -f coverage.out
