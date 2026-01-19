package tests

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
}

func TestAddNewCatalogItem(t *testing.T) {
	db, repo := setupTest(t)

	newItem := &pb.CatalogItem{
		ItemId:            "item789",
		Description:       "New Test Item",
		QuantityAvailable: 20,
		Price:             29.99,
	}
	err := repo.AddCatalogItem(newItem)
	if err != nil {
		t.Errorf("Failed to add new catalog item: %v", err)
	}

	// Verify the item was added in the database
	var item domain.CatalogItem
	err = db.Where("item_id = ?", "item789").First(&item).Error
	if err != nil {
		t.Errorf("Error retrieving added item: %v", err)
	}
	if item.Description != "New Test Item" {
		t.Errorf("Added item description mismatch: got %v, want %v", item.Description, "New Test Item")
	}
	if item.QuantityAvailable != 20 {
		t.Errorf("Added item quantity mismatch: got %v, want %v", item.QuantityAvailable, 20)
	}
	if item.Price != 29.99 {
		t.Errorf("Added item price mismatch: got %v, want %v", item.Price, 29.99)
	}
}

func TestAddExistingCatalogItem(t *testing.T) {
	db, repo := setupTest(t)

	existingItem := &pb.CatalogItem{
		ItemId:            "item123",
		Description:       "Updated Default Item",
		QuantityAvailable: 15,
		Price:             89.99,
	}
	err := repo.AddCatalogItem(existingItem)
	if err == nil {
		t.Errorf("Expected error when adding existing catalog item, but got none")
	}

	// Verify that the database wasn't altered
	var item domain.CatalogItem
	err = db.Where("item_id = ?", "item123").First(&item).Error
	if err != nil {
		t.Errorf("Error retrieving updated item: %v", err)
	}
	if item.Description != "Default Item" {
		t.Errorf("Existing item description should not be updated: got %v, want %v", item.Description, "Default Item")
	}
	if item.QuantityAvailable != 10 {
		t.Errorf("Existing item quantity should not be updated: got %v, want %v", item.QuantityAvailable, 10)
	}
	if item.Price != 99.99 {
		t.Errorf("Existing item price should not be updated: got %v, want %v", item.Price, 99.99)
	}
}

func TestAddCatalogItemInvalidID(t *testing.T) {
	_, repo := setupTest(t)

	invalidItem := &pb.CatalogItem{
		ItemId:            "",
		Description:       "Invalid Item",
		QuantityAvailable: 10,
		Price:             19.99,
	}

	err := repo.AddCatalogItem(invalidItem)
	if err == nil {
		t.Errorf("Expected error: %v, but got none", err)
	}
}

func TestAddCatalogItemInvalidDescription(t *testing.T) {
	_, repo := setupTest(t)

	invalidItem := &pb.CatalogItem{
		ItemId:            "item999",
		Description:       "",
		QuantityAvailable: 10,
		Price:             19.99,
	}

	err := repo.AddCatalogItem(invalidItem)
	if err == nil {
		t.Errorf("Expected error: %v, but got none", err)
	}
}

func TestAddCatalogItemInvalidQuantity(t *testing.T) {
	_, repo := setupTest(t)

	invalidItem := &pb.CatalogItem{
		ItemId:            "item999",
		Description:       "Invalid Quantity Item",
		QuantityAvailable: 0,
		Price:             19.99,
	}

	err := repo.AddCatalogItem(invalidItem)
	if err == nil {
		t.Errorf("Expected error: %v, but got none", err)
	}
}

func TestAddCatalogItemInvalidPrice(t *testing.T) {
	_, repo := setupTest(t)

	invalidItem := &pb.CatalogItem{
		ItemId:            "item999",
		Description:       "Invalid Price Item",
		QuantityAvailable: 10,
		Price:             -5.00,
	}

	err := repo.AddCatalogItem(invalidItem)
	if err == nil {
		t.Errorf("Expected error: %v, but got none", err)
	}
}

func TestRemoveExistingCatalogItem(t *testing.T) {
	db, repo := setupTest(t)

	err := repo.RemoveCatalogItem("item456")
	if err != nil {
		t.Errorf("Failed to remove existing catalog item: %v", err)
	}

	// Verify the item was removed from the database
	var item domain.CatalogItem
	err = db.Where("item_id = ?", "item456").First(&item).Error
	if err == nil {
		t.Errorf("Item was not removed from the database")
	}
}

func TestRemoveNonExistingCatalogItem(t *testing.T) {
	_, repo := setupTest(t)

	err := repo.RemoveCatalogItem("nonexistent_item")
	if err == nil {
		t.Errorf("Expected error: %v but got none", err)
	}
}

func TestRemoveCatalogItemInvalidID(t *testing.T) {
	_, repo := setupTest(t)

	err := repo.RemoveCatalogItem("")
	if err == nil {
		t.Errorf("Expected error: %v, but got none", err)
	}
}

func TestGetExistingCatalogItem(t *testing.T) {
	_, repo := setupTest(t)

	item, err := repo.GetCatalogItem("item123")
	if err != nil {
		t.Errorf("Failed to get existing catalog item: %v", err)
	}
	if item.ItemId != "item123" {
		t.Errorf("Retrieved item ID mismatch: got %v, want %v", item.ItemId, "item123")
	}
	if item.Description != "Default Item" {
		t.Errorf("Retrieved item description mismatch: got %v, want %v", item.Description, "Default Item")
	}
	if item.QuantityAvailable != 10 {
		t.Errorf("Retrieved item quantity mismatch: got %v, want %v", item.QuantityAvailable, 10)
	}
	if item.Price != 99.99 {
		t.Errorf("Retrieved item price mismatch: got %v, want %v", item.Price, 99.99)
	}
}

func TestGetNonExistingCatalogItem(t *testing.T) {
	_, repo := setupTest(t)
	_, err := repo.GetCatalogItem("nonexistent_item")
	if err == nil {
		t.Errorf("Expected error: %v but got none", err)
	}
}

func TestGetCatalogItemInvalidID(t *testing.T) {
	_, repo := setupTest(t)
	_, err := repo.GetCatalogItem("")
	if err == nil {
		t.Errorf("Expected error: %v, but got none", err)
	}
}

func TestUpdateQuantityAvailableValid(t *testing.T) {
	db, repo := setupTest(t)

	err := repo.UpdateQuantityAvailable("item123", 25)
	if err != nil {
		t.Errorf("Failed to update quantity available: %v", err)
	}

	var item domain.CatalogItem
	err = db.Where("item_id = ?", "item123").First(&item).Error
	if err != nil {
		t.Errorf("Failed to retrieve updated item: %v", err)
	}
	if item.QuantityAvailable != 25 {
		t.Errorf("Quantity available not updated correctly: got %v, want %v", item.QuantityAvailable, 25)
	}
}

func TestUpdateQuantityAvailableInvalid(t *testing.T) {
	db, repo := setupTest(t)

	err := repo.UpdateQuantityAvailable("item123", 0)
	if err == nil {
		t.Errorf("Expected error: %v, but got none", err)
	}

	var item domain.CatalogItem
	err = db.Where("item_id = ?", "item123").First(&item).Error
	if err != nil {
		t.Errorf("Failed to retrieve item after invalid update attempt: %v", err)
	}
	if item.QuantityAvailable != 10 {
		t.Errorf("Quantity available should not be updated on invalid input: got %v, want %v", item.QuantityAvailable, 10)
	}
}

func TestUpdateQuantityAvailableInvalidID(t *testing.T) {
	_, repo := setupTest(t)

	err := repo.UpdateQuantityAvailable("", 15)
	if err == nil {
		t.Errorf("Expected error: %v, but got none", err)
	}
}

func TestUpdateQuantityAvailableNonExistingItem(t *testing.T) {
	_, repo := setupTest(t)

	err := repo.UpdateQuantityAvailable("nonexistent_item", 15)
	if err == nil {
		t.Errorf("Expected error: %v but got none", err)
	}
}

func TestUpdatePriceValid(t *testing.T) {
	db, repo := setupTest(t)

	err := repo.UpdatePrice("item123", 79.99)
	if err != nil {
		t.Errorf("Failed to update price: %v", err)
	}

	var item domain.CatalogItem
	err = db.Where("item_id = ?", "item123").First(&item).Error
	if err != nil {
		t.Errorf("Failed to retrieve updated item: %v", err)
	}
	if item.Price != 79.99 {
		t.Errorf("Price not updated correctly: got %v, want %v", item.Price, 79.99)
	}
}

func TestUpdatePriceInvalid(t *testing.T) {
	db, repo := setupTest(t)

	err := repo.UpdatePrice("item123", -10.00)
	if err == nil {
		t.Errorf("Expected error: %v, but got none", err)
	}

	var item domain.CatalogItem
	err = db.Where("item_id = ?", "item123").First(&item).Error
	if err != nil {
		t.Errorf("Failed to retrieve item after invalid update attempt: %v", err)
	}
	if item.Price != 99.99 {
		t.Errorf("Price should not be updated on invalid input: got %v, want %v", item.Price, 99.99)
	}
}

func TestUpdatePriceInvalidID(t *testing.T) {
	_, repo := setupTest(t)

	err := repo.UpdatePrice("", 49.99)
	if err == nil {
		t.Errorf("Expected error: %v, but got none", err)
	}
}

func TestUpdatePriceNonExistingItem(t *testing.T) {
	_, repo := setupTest(t)

	err := repo.UpdatePrice("nonexistent_item", 49.99)
	if err == nil {
		t.Errorf("Expected error: %v but got none", err)
	}
}

func TestListCatalogItems(t *testing.T) {
	_, repo := setupTest(t)

	items, err := repo.ListCatalogItems()
	if err != nil {
		t.Errorf("Failed to list catalog items: %v", err)
	}
	if len(items) != 2 {
		t.Errorf("Expected 2 catalog items, got %v", len(items))
	}

	expectedItems := map[string]*pb.CatalogItem{
		"item123": {
			ItemId:            "item123",
			Description:       "Default Item",
			QuantityAvailable: 10,
			Price:             99.99,
		},
		"item456": {
			ItemId:            "item456",
			Description:       "Another Item",
			QuantityAvailable: 5,
			Price:             49.99,
		},
	}
	for _, item := range items {
		expectedItem, exists := expectedItems[item.ItemId]
		if !exists {
			t.Errorf("Unexpected item ID: %v", item.ItemId)
			continue
		}
		if item.Description != expectedItem.Description {
			t.Errorf("Item description mismatch for %v: got %v, want %v", item.ItemId, item.Description, expectedItem.Description)
		}
		if item.QuantityAvailable != expectedItem.QuantityAvailable {
			t.Errorf("Item quantity mismatch for %v: got %v, want %v", item.ItemId, item.QuantityAvailable, expectedItem.QuantityAvailable)
		}
		if item.Price != expectedItem.Price {
			t.Errorf("Item price mismatch for %v: got %v, want %v", item.ItemId, item.Price, expectedItem.Price)
		}
	}
}
