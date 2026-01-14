# CONFIG
APP_NAME=ecommerce
BIN_DIR=bin

GO=go
GOFLAGS=-v
ENV?=dev

PROTOC=protoc

GO_OUT=.
GO_OPT=paths=source_relative

PROTO_DIR=ecommerce/proto
PROTO_FILES=$(wildcard $(PROTO_DIR)/*.proto)

ifeq ($(OS),Windows_NT)
	CLEAN_PROTO = powershell -Command "Remove-Item -Force -ErrorAction SilentlyContinue $(PROTO_DIR)\*.pb.go"
else
	CLEAN_PROTO = rm -f $(PROTO_DIR)/*.pb.go
endif

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

# PROTOBUF
.PHONY: proto
proto:
	$(PROTOC) \
		--go_out=$(GO_OUT) --go_opt=$(GO_OPT) \
		--go-grpc_out=$(GO_OUT) --go-grpc_opt=$(GO_OPT) \
		$(PROTO_FILES)

.PHONY: clean-proto
clean-proto:
	$(CLEAN_PROTO)

# TESTING
.PHONY: test
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
