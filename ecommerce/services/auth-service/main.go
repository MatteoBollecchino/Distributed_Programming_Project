package main

// We need to import the driver for SQLite
// to use GORM

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	pb "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/auth"
	"github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/services/auth-service/internal"
	"github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/services/auth-service/internal/repository"
)

var port = "8081"

func main() {

	// Initialize database connection with GORM
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// Start gRPC server
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	// Initialize repository and create default users/admins
	authRepo := repository.NewAuthRepository(db)
	if err := authRepo.CreateDefaultUsersAdmins(); err != nil {
		log.Fatalf("Internal errors while creating default users: %v", err)
	}

	// Initialize AuthServer
	authServer := internal.NewAuthServer(authRepo)

	// Register gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterAuthenticationServiceServer(grpcServer, authServer)

	/*
		// gRPC Health Server
		healthServer := health.NewServer()
		healthServer.SetServingStatus("auth.AuthenticationService", healthpb.HealthCheckResponse_SERVING)
		healthpb.RegisterHealthServer(grpcServer, healthServer)
	*/

	log.Printf("Auth service listening on port %s", port)

	// Start serving
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC user server failed: %v", err)
	}
}
