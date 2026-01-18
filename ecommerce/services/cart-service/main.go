package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	pb "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/cart"
	"github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/services/cart-service/internal"
	"github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/services/cart-service/internal/domain"
	"github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/services/cart-service/internal/repository"
)

var port = "8082"

func main() {

	// Initialize database connection with GORM
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// Migrate the schema
	if err := db.AutoMigrate(&domain.Cart{}, &domain.CartItem{}); err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}

	// Start gRPC server
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	// Initialize repository
	cartRepo := repository.NewCartServiceRepository(db)

	// Initialize CartServer
	cartServer := internal.NewCartServer(cartRepo)

	// Register gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterCartServiceServer(grpcServer, cartServer)

	log.Printf("Cart service listening on port %s", port)

	// Start serving
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC cart service failed: %v", err)
	}
}
