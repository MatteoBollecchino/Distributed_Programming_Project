# ==========================
# SETUP
# ==========================
OS := $(shell go env GOOS)

ifeq ($(OS),windows)
	RM = powershell -Command "Remove-Item -Recurse -Force -ErrorAction SilentlyContinue"
	FIND_PROTO = powershell -Command "Get-ChildItem proto -Recurse -Include *.pb.go,*_grpc.pb.go | Remove-Item -Force"
	EXE = .exe
else
	RM = rm -f
	FIND_PROTO = find proto -type f \( -name "*.pb.go" -o -name "*_grpc.pb.go" \) -delete
	EXE =
endif


# ==========================
# VARIABLES
# ==========================
SERVICES_DIR=services

SERVICES=auth-service \
         cart-service \
         catalog-service \
         order-service \
         payment-service

PROTOC=protoc

PROTO_SERVICES=auth cart catalog order payment

GO=go

# ==========================
# DEFAULT
# ==========================
.PHONY: all
all: proto build

# ==========================
# PROTOBUF
# ==========================
.PHONY: proto
proto:
	@echo "Generating protobuf files..."
	$(PROTOC) --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/auth/auth.proto
	$(PROTOC) --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/cart/cart.proto
	$(PROTOC) --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/catalog/catalog.proto
	$(PROTOC) --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/order/order.proto
	$(PROTOC) --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/payment/payment.proto
	@echo "Protobuf generated"

# ==========================
# CLEAN PROTO GENERATED FILES
# ==========================
.PHONY: clean-proto
clean-proto:
	@echo "Removing generated protobuf files..."
	@$(FIND_PROTO)
	@echo "Protobuf files removed"

# ==========================
# BUILD
# ==========================
.PHONY: build
build:
	@echo "Building services..."
	@for svc in $(SERVICES); do \
		cd $(SERVICES_DIR)/$$svc && $(GO) build -o $$svc ; \
	done
	@echo "Build completed"

# ==========================
# RUN SINGLE SERVICES
# ==========================
.PHONY: run-auth
run-auth:
	cd $(SERVICES_DIR)/auth-service && $(GO) run main.go

.PHONY: run-cart
run-cart:
	cd $(SERVICES_DIR)/cart-service && $(GO) run main.go

.PHONY: run-catalog
run-catalog:
	cd $(SERVICES_DIR)/catalog-service && $(GO) run main.go

.PHONY: run-order
run-order:
	cd $(SERVICES_DIR)/order-service && $(GO) run main.go

.PHONY: run-payment
run-payment:
	cd $(SERVICES_DIR)/payment-service && $(GO) run main.go

# ==========================
# RUN ALL SERVICES
# ==========================
.PHONY: run-all
run-all:
	go run scripts/run-all.go

# ==========================
# TEST
# ==========================
.PHONY: test
test:
	@echo "Running tests..."
	@for svc in $(SERVICES); do \
		cd $(SERVICES_DIR)/$$svc && $(GO) test ./... ; \
	done
	@echo "Tests completed"

# ==========================
# CLEAN
# ==========================
.PHONY: clean
clean:
	@echo "Cleaning binaries..."
	@for svc in $(SERVICES); do \
		$(RM) $(SERVICES_DIR)/$$svc/$$svc$(EXE); \
	done
	@echo "Binaries removed"
