package main

// We need to import the driver for SQLite
// to use GORM

import (
	"github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/services/auth-service/internal"
)

func main() {
	authServer := internal.NewAuthServer(nil)
	_ = authServer
	/*
		db, _ := database.New("POSTGRES_DSN")

		userRepo := repository.NewUserRepository(db)
		loginUC := auth.NewLoginUseCase(userRepo)
		handler := http.NewHandler(loginUC)

		router := http.NewRouter(handler)

		log.Println("Auth service running on :8081")
		http.ListenAndServe(":8081", router)
	*/
}
