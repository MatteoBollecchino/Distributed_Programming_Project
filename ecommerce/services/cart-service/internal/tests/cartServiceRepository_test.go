package tests

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/services/cart-service/internal/domain"
	"github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/services/cart-service/internal/repository"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect database: %v", err)
	}

	err = db.AutoMigrate(&domain.Cart{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}
	return db
}

func setupDefaultCarts(t *testing.T, db *gorm.DB, repo *repository.CartServiceRepository) {

}

func setupTest(t *testing.T) (*gorm.DB, *repository.CartServiceRepository) {
	db := setupTestDB(t)
	repo := repository.NewCartServiceRepository(db)

	setupDefaultCarts(t, db, repo)

	return db, repo
}
