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
TESTS_DIR=internal/tests
PROTOC=protoc
GO=go

# ==========================
# DEFAULT
# ==========================
.PHONY: all proto build clean clean-proto test stop-all run-all run-all-tabs
all: proto build

# ==========================
# RUN ALL SERVICES
# ==========================
run-all:
	@echo "Starting all services in background..."
	cd $(SERVICES_DIR)/auth-service && go run main.go &
	cd $(SERVICES_DIR)/cart-service && go run main.go &
	cd $(SERVICES_DIR)/catalog-service && go run main.go &
	cd $(SERVICES_DIR)/order-service && go run main.go &
	cd $(SERVICES_DIR)/payment-service && go run main.go &
	@echo "All services are running."

# Logica condizionale spostata FUORI dal target per non rompere il parsing di Make
run-all-tabs:
ifeq ($(OS),windows)
	@echo "Starting services in Windows Terminal..."
	@start "Auth Service" cmd /k "cd $(CURDIR)\$(SERVICES_DIR)\auth-service && $(GO) run main.go"
	@start "Cart Service" cmd /k "cd $(CURDIR)\$(SERVICES_DIR)\cart-service && $(GO) run main.go"
	@start "Catalog Service" cmd /k "cd $(CURDIR)\$(SERVICES_DIR)\catalog-service && $(GO) run main.go"
	@start "Order Service" cmd /k "cd $(CURDIR)\$(SERVICES_DIR)\order-service && $(GO) run main.go"
	@start "Payment Service" cmd /k "cd $(CURDIR)\$(SERVICES_DIR)\payment-service && $(GO) run main.go"
else
ifeq ($(shell uname),Linux)
	@echo "Starting services in GNOME Terminal tabs..."
	@gnome-terminal \
		--tab --title="Auth" -- bash -c "cd $(CURDIR)/$(SERVICES_DIR)/auth-service && $(GO) run main.go; exec bash" \
		--tab --title="Cart" -- bash -c "cd $(CURDIR)/$(SERVICES_DIR)/cart-service && $(GO) run main.go; exec bash" \
		--tab --title="Catalog" -- bash -c "cd $(CURDIR)/$(SERVICES_DIR)/catalog-service && $(GO) run main.go; exec bash" \
		--tab --title="Order" -- bash -c "cd $(CURDIR)/$(SERVICES_DIR)/order-service && $(GO) run main.go; exec bash" \
		--tab --title="Payment" -- bash -c "cd $(CURDIR)/$(SERVICES_DIR)/payment-service && $(GO) run main.go; exec bash"
else
	@echo "You're on MACOS"
endif
endif

# ==========================
# PROTOBUF
# ==========================
proto:
	@echo "Generating protobuf files..."
	$(PROTOC) --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/auth/auth.proto
	$(PROTOC) --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/cart/cart.proto
	$(PROTOC) --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/catalog/catalog.proto
	$(PROTOC) --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/order/order.proto
	$(PROTOC) --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/payment/payment.proto
	@echo "Protobuf generated"

clean-proto:
	@echo "Removing generated protobuf files..."
	@$(FIND_PROTO)
	@echo "Protobuf files removed"

# ==========================
# BUILD
# ==========================
build:
	@echo "Building services..."
	cd $(SERVICES_DIR)/auth-service && $(GO) build -o auth-service
	cd $(SERVICES_DIR)/cart-service && $(GO) build -o cart-service
	cd $(SERVICES_DIR)/catalog-service && $(GO) build -o catalog-service
	cd $(SERVICES_DIR)/order-service && $(GO) build -o order-service
	cd $(SERVICES_DIR)/payment-service && $(GO) build -o payment-service
	@echo "Build completed"

# ==========================
# RUN SINGLE SERVICES
# ==========================
run-auth:
	@echo "Starting Authentification Service..."
	cd $(SERVICES_DIR)/auth-service && $(GO) run main.go

run-cart:
	@echo "Starting Cart Service..."
	cd $(SERVICES_DIR)/cart-service && $(GO) run main.go

run-catalog:
	@echo "Starting Catalog Service..."
	cd $(SERVICES_DIR)/catalog-service && $(GO) run main.go

run-order:
	@echo "Starting Order Service..."
	cd $(SERVICES_DIR)/order-service && $(GO) run main.go

run-payment:
	@echo "Starting Payment Service..."
	cd $(SERVICES_DIR)/payment-service && $(GO) run main.go

# ==========================
# STOP ALL SERVICES
# ==========================
stop-all:
	@echo "Stopping all services..."
ifeq ($(OS),windows)
	@taskkill /F /IM main.exe /T 2>nul || true
	@taskkill /F /IM auth-service.exe /T 2>nul || true
	@taskkill /F /IM cart-service.exe /T 2>nul || true
	@taskkill /F /IM catalog-service.exe /T 2>nul || true
	@taskkill /F /IM order-service.exe /T 2>nul || true
	@taskkill /F /IM payment-service.exe /T 2>nul || true
else
	@-pkill -f "go run main.go"
	@-fuser -k 8081/tcp 2>/dev/null || true
	@-fuser -k 8082/tcp 2>/dev/null || true
	@-fuser -k 8083/tcp 2>/dev/null || true
	@-fuser -k 8084/tcp 2>/dev/null || true
	@-fuser -k 8085/tcp 2>/dev/null || true
endif
	@echo "Services stopped."

# ==========================
# TEST
# ==========================
test:
	@echo "Running tests..."
	cd $(SERVICES_DIR)/auth-service/$(TESTS_DIR) && $(GO) test ./...
	cd $(SERVICES_DIR)/cart-service/$(TESTS_DIR) && $(GO) test ./...
	cd $(SERVICES_DIR)/catalog-service/$(TESTS_DIR) && $(GO) test ./...
	cd $(SERVICES_DIR)/order-service/$(TESTS_DIR) && $(GO) test ./...
	cd $(SERVICES_DIR)/payment-service/$(TESTS_DIR) && $(GO) test ./...
	@echo "Tests completed"

# ==========================
# CLEAN
# ==========================
clean:
	@echo "Cleaning binaries..."
	$(RM) $(SERVICES_DIR)/auth-service/auth-service$(EXE)
	$(RM) $(SERVICES_DIR)/cart-service/cart-service$(EXE)
	$(RM) $(SERVICES_DIR)/catalog-service/catalog-service$(EXE)
	$(RM) $(SERVICES_DIR)/order-service/order-service$(EXE)
	$(RM) $(SERVICES_DIR)/payment-service/payment-service$(EXE)
	@echo "Binaries removed"