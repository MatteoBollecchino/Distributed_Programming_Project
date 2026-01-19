package tests

/*
import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	pb "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/catalog"
	"github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/services/catalog-service/internal/domain"
	"github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/services/catalog-service/internal/repository"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect database: %v", err)
	}

	err = db.AutoMigrate(&domain.CatalogItem{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}
	return db
}

func setupDefaultCatalogItems(db *gorm.DB) {
	defaultItem1 := &pb.CatalogItem{
		ItemId:            "item123",
		Description:       "Default Item",
		QuantityAvailable: 10,
		Price:             99.99,
	}

	defaultItem2 := &pb.CatalogItem{
		ItemId:            "item456",
		Description:       "Another Item",
		QuantityAvailable: 5,
		Price:             49.99,
	}

	db.Create(defaultItem1)
	db.Create(defaultItem2)
}

func setupTest(t *testing.T) (*gorm.DB, *repository.CatalogServiceRepository) {
	db := setupTestDB(t)
	repo := repository.NewCatalogServiceRepository(db)

	setupDefaultCatalogItems(db)

	return db, repo
}*/
