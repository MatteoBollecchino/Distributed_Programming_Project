# Distributed_Programming_Project

The idea underlying the Distributed Programming project is the development of an e-commerce
platform dedicated to the sale of products related to the fantasy and science fiction world, such as
books, manga, collectible items, etc.
The main objective of the project is to design and implement a realistic distributed system capable
of concretely applying the concepts covered during the course, including inter-process communication,
concurrency management, separation of concerns, and architectural scalability.

## Prerequisites

- Go 1.21 or higher

## Makefile instructions 

make run            # Start
make build          # Build binary
make test           # Run tests
make test-cover     # Run tests with coverage
make clean          # Clean build artifacts

## Project Flow

Browser
  ↓
Web Server (HTML + sessioni)
  ↓ REST
Auth / Catalog / Cart Service
  ↓ gRPC
Order Service ↔ Catalog Service
  ↓
Database (GORM)


## Project Structure

dp-ecommerce/
├── services/
│   ├── auth-service/
│   ├── catalog-service/
│   ├── cart-service/
│   ├── order-service/
│   └── wishlist-service/
├── web/
│   ├── templates/
│   └── server/
├── proto/
│   ├── catalog.proto
│   ├── order.proto
│   └── cart.proto
├── shared/
│   ├── logger/
│   ├── middleware/
│   └── auth/
├── Makefile
└── README.md



