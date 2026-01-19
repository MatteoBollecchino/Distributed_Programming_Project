package main

/*
import (
	"log"
	"net"

	"google.golang.org/grpc"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	pb "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/catalog"
	"github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/services/catalog-service/internal"
	"github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/services/catalog-service/internal/domain"
	"github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/services/catalog-service/internal/repository"
)

var port = "8083"

func main() {

	// Initialize database connection with GORM
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// Migrate the schema
	if err := db.AutoMigrate(&domain.CatalogItem{}); err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}

	// Start gRPC server
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	// Initialize repository
	catalogRepo := repository.NewCatalogServiceRepository(db)

	// Initialize CatalogServer
	catalogServer := internal.NewCatalogServer(catalogRepo)

	// Register gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterCatalogServiceServer(grpcServer, catalogServer)

	log.Printf("Catalog service listening on port %s", port)

	// Start serving
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC catalog service failed: %v", err)
	}
}*/
