package main

/*
import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/sqlite" // We need to import the driver for SQLite
	"gorm.io/gorm"          // to use GORM
)
*/

/*
import (
	"log"
	"net/http"

	"services/auth-service/internal/application/auth"
	"services/auth-service/internal/infrastructure/database"
	"services/auth-service/internal/infrastructure/http"
	"services/auth-service/internal/infrastructure/repository"
)

func main() {
	db, _ := database.New("POSTGRES_DSN")

	userRepo := repository.NewUserRepository(db)
	loginUC := auth.NewLoginUseCase(userRepo)
	handler := http.NewHandler(loginUC)

	router := http.NewRouter(handler)

	log.Println("Auth service running on :8081")
	http.ListenAndServe(":8081", router)
}
*/
