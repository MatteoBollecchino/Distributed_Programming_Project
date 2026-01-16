package tests

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	pb "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/cart"
	"github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/services/cart-service/internal/domain"
	"github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/services/cart-service/internal/repository"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect database: %v", err)
	}

	err = db.AutoMigrate(&domain.Cart{}, &domain.CartItem{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}
	return db
}

func setupDefaultCarts(t *testing.T, db *gorm.DB, repo *repository.CartServiceRepository) {

	cart1 := &domain.Cart{
		Username: "user1",
		Items: []domain.CartItem{
			{ItemID: "item1", Quantity: 2, Price: 10.0},
			{ItemID: "item2", Quantity: 1, Price: 20.0},
		},
	}

	cart2 := &domain.Cart{
		Username: "user2",
		Items: []domain.CartItem{
			{ItemID: "item3", Quantity: 5, Price: 5.0},
			{ItemID: "item4", Quantity: 2, Price: 15.0},
		},
	}

	err := db.Create(cart1).Error
	if err != nil {
		t.Fatalf("Failed to create cart1: %v", err)
	}

	err = db.Create(cart2).Error
	if err != nil {
		t.Fatalf("Failed to create cart2: %v", err)
	}
}

func setupTest(t *testing.T) (*gorm.DB, *repository.CartServiceRepository) {
	db := setupTestDB(t)
	repo := repository.NewCartServiceRepository(db)

	setupDefaultCarts(t, db, repo)

	return db, repo
}

func TestAddNewItemToCart(t *testing.T) {
	db, repo := setupTest(t)

	// Test adding an item to an existing cart
	cartItem1 := &pb.CartItem{ItemId: "item3", Quantity: 1, Price: 30.0}
	err := repo.AddItemToCart("user1", cartItem1)
	if err != nil {
		t.Errorf("Failed to add item to existing cart: %v", err)
	}

	// Test if item was added correctly in the database
	var cart domain.Cart
	err = db.Where("username = ?", "user1").First(&cart).Error
	if err != nil {
		t.Errorf("Failed to retrieve cart from database: %v", err)
	}
	if len(cart.Items) != 3 {
		t.Errorf("Expected 3 items in cart, got %d", len(cart.Items))
	}
}

// DA CONTROLLARE
func TestAddExistingItemToCart(t *testing.T) {
	_, repo := setupTest(t)

	// Test adding an existing item to the cart (should update quantity)
	cartItem := &pb.CartItem{ItemId: "item1", Quantity: 3, Price: 10.0}
	err := repo.AddItemToCart("user1", cartItem)
	if err != nil {
		t.Errorf("Failed to add existing item to cart: %v", err)
	}
}

// Incorrect adding
