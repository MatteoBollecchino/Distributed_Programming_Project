package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	pb "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/payment"
	"github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/services/payment-service/internal"
	"github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/services/payment-service/internal/domain"
	"github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/services/payment-service/internal/repository"
)

var port = "8085"

func main() {

	// Initialize database connection with GORM
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// Migrate the schema
	if err := db.AutoMigrate(&domain.Payment{}); err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}

	// Start gRPC server
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	// Initialize repository
	paymentRepo := repository.NewPaymentServiceRepository(db)

	// Initialize PaymentServer
	paymentServer := internal.NewPaymentServer(paymentRepo)

	// Register gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterPaymentServiceServer(grpcServer, paymentServer)

	log.Printf("Payment service listening on port %s", port)

	// Start serving
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC payment service failed: %v", err)
	}
}
