# Distributed_Programming_Project

The idea underlying the Distributed Programming project is the development of an e-commerce
platform dedicated to the sale of products related to the fantasy and science fiction world, such as
books, manga, collectible items, etc.
The main objective of the project is to design and implement a realistic distributed system capable
of concretely applying the concepts covered during the course, including inter-process communication,
concurrency management, separation of concerns, and architectural scalability.

## Prerequisites

- Go 1.21 or higher
- protoc compiler

## Makefile instructions 

make run-all        # Start
make proto          # Compiles with protoc
make clean-proto    # Removes files generated with protoc
make build          # Build binary
make test           # Run tests
make clean          # Clean build artifacts

## Project Flow

Browser
  ↓
Web Server (HTML + sessions)
  ↓ REST
Auth / Catalog / Cart Service
  ↓ gRPC
Order Service ↔ Catalog Service
  ↓
Database (GORM)

# Request Flow

GET /catalog
  ↓
RequireAuth middleware
  ↓
CatalogHandler
  ↓
CatalogClient.ListProducts()
  ↓
Catalog Service
  ↓
ProductDTO[]
  ↓
ProductViewModel[]
  ↓
catalog.html


## Project Structure

METTERE LO SCHEMA CHE SI OTTIENE DAL COMANDO "tree -F ecommerce"

ecommerce/
├── services/
│   ├── auth-service/
│   ├── catalog-service/
│   ├── cart-service/
│   ├── order-service/
│   └── payment-service/
├── web/
│   ├── templates/
│   └── server/
├── proto/
│   ├── catalog.proto
│   ├── order.proto
│   └── cart.proto
├── Makefile
└── README.md

*-service/
├── internal/
│   ├── domain/
│   │   ├── entity.go
│   │   ├── interface.go
|   | 
│   ├── repository/
│   │   └── product_repository.go
│   └── tests/
├── go.mod
└── Makefile

web/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── handlers/
│   │   ├── auth.go
│   │   ├── catalog.go
│   │   ├── cart.go
│   │   ├── order.go
│   │   └── admin.go
│   ├── clients/
│   │   ├── auth_client.go
│   │   ├── catalog_client.go
│   │   ├── cart_client.go
│   │   └── order_client.go
│   ├── session/
│   │   └── secure_cookie.go
│   ├── middleware/
│   │   ├── auth.go
│   │   └── admin.go
│   └── viewmodels/
│       ├── product_vm.go
│       └── order_vm.go
├── templates/
│   ├── login.html
│   ├── catalog.html
│   ├── cart.html
│   └── orders.html





