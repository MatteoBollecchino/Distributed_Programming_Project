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
			{ItemID: "item1", CartUsername: "user1", Quantity: 2, Price: 10.0},
			{ItemID: "item2", CartUsername: "user1", Quantity: 1, Price: 20.0},
		},
	}

	cart2 := &domain.Cart{
		Username: "user2",
		Items: []domain.CartItem{
			{ItemID: "item3", CartUsername: "user2", Quantity: 5, Price: 5.0},
			{ItemID: "item4", CartUsername: "user2", Quantity: 2, Price: 15.0},
		},
	}

	if err := db.Session(&gorm.Session{FullSaveAssociations: true}).Create(cart1).Error; err != nil {
		t.Fatalf("Failed to create cart1: %v", err)
	}

	if err := db.Session(&gorm.Session{FullSaveAssociations: true}).Create(cart2).Error; err != nil {
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
	err = db.Preload("Items").Where("username = ?", "user1").First(&cart).Error
	if err != nil {
		t.Errorf("Failed to retrieve cart from database: %v", err)
	}
	if len(cart.Items) != 3 {
		t.Errorf("Expected 3 items in cart, got %d", len(cart.Items))
	}

}

func TestAddExistingItemToCart(t *testing.T) {
	db, repo := setupTest(t)

	// Verify initial quantity of item1 in user1's cart
	var initialCart domain.Cart
	err := db.Preload("Items").Where("username = ?", "user1").First(&initialCart).Error
	if err != nil {
		t.Errorf("Failed to retrieve cart from database: %v", err)
	}

	// Test adding an existing item to the cart (should update quantity)
	cartItem := &pb.CartItem{ItemId: "item1", Quantity: 3, Price: 10.0}
	err = repo.AddItemToCart("user1", cartItem)
	if err != nil {
		t.Errorf("Failed to add existing item to cart: %v", err)
	}

	// Verify if the quantity was updated correctly
	var cart domain.Cart
	err = db.Preload("Items").Where("username = ?", "user1").First(&cart).Error
	if err != nil {
		t.Errorf("Failed to retrieve cart from database: %v", err)
	}
	for _, item := range cart.Items {
		if item.ItemID == "item1" {
			if item.Quantity != 5 {
				t.Errorf("Expected quantity 5 for item1, got %d", item.Quantity)
			}
		}
	}
}

func TestAddNewItemToNewCart(t *testing.T) {
	db, repo := setupTest(t)

	// Test adding an item to a new cart (cart does not exist yet)
	cartItem := &pb.CartItem{ItemId: "item1", Quantity: 2, Price: 50.0}
	err := repo.AddItemToCart("newuser", cartItem)
	if err != nil {
		t.Errorf("Failed to add item to new cart: %v", err)
	}

	// Verify if the new cart was created correctly
	var cart domain.Cart
	err = db.Preload("Items").Where("username = ?", "newuser").First(&cart).Error
	if err != nil {
		t.Errorf("Failed to retrieve new cart from database: %v", err)
	}
	if len(cart.Items) != 1 {
		t.Errorf("Expected 1 item in new cart, got %d", len(cart.Items))
	}
}

func TestAddItemToCartWithEmptyID(t *testing.T) {
	db, repo := setupTest(t)

	// Test adding an item with nil ID
	cartItem := &pb.CartItem{ItemId: "", Quantity: 2, Price: 10.0}
	err := repo.AddItemToCart("user1", cartItem)
	if err == nil {
		t.Errorf("Expected error when adding item with nil ID, got nil")
	}

	// Verify that the cart remains unchanged
	var cart domain.Cart
	err = db.Preload("Items").Where("username = ?", "user1").First(&cart).Error
	if err != nil {
		t.Errorf("Failed to retrieve cart from database: %v", err)
	}
	if len(cart.Items) != 2 {
		t.Errorf("Expected 2 items in cart after failed addition, got %d", len(cart.Items))
	}
}

func TestAddItemToCartWithEmptyCartUsername(t *testing.T) {
	db, repo := setupTest(t)

	// Test adding an item with nil ID
	cartItem := &pb.CartItem{ItemId: "item6", Quantity: 2, Price: 10.0}
	err := repo.AddItemToCart("", cartItem)
	if err == nil {
		t.Errorf("Expected error when adding item with nil ID, got nil")
	}

	// Verify that the cart remains unchanged
	var cart domain.Cart
	err = db.Preload("Items").Where("username = ?", "").First(&cart).Error
	if err == nil {
		t.Errorf("Failed to retrieve cart from database: %v", err)
	}
	if len(cart.Items) != 0 {
		t.Errorf("Expected 0 items in cart after failed addition, got %d", len(cart.Items))
	}
}

func TestAddItemToCartWithZeroQuantity(t *testing.T) {
	db, repo := setupTest(t)

	// Test adding an item with zero quantity
	cartItem := &pb.CartItem{ItemId: "itemX", Quantity: 0, Price: 10.0}
	err := repo.AddItemToCart("user1", cartItem)
	if err == nil {
		t.Errorf("Expected error when adding item with zero quantity, got nil")
	}

	// Verify that the cart remains unchanged
	var cart domain.Cart
	err = db.Preload("Items").Where("username = ?", "user1").First(&cart).Error
	if err != nil {
		t.Errorf("Failed to retrieve cart from database: %v", err)
	}
	if len(cart.Items) != 2 {
		t.Errorf("Expected 2 items in cart after failed addition, got %d", len(cart.Items))
	}
}

func TestAddItemToCartWithZeroPrice(t *testing.T) {
	db, repo := setupTest(t)

	// Test adding an item with zero price
	cartItem := &pb.CartItem{ItemId: "itemX", Quantity: 1, Price: 0.0}
	err := repo.AddItemToCart("user1", cartItem)
	if err != nil {
		t.Errorf("Expected error when adding item with zero price, got nil")
	}

	// Verify that the cart changes
	var cart domain.Cart
	err = db.Preload("Items").Where("username = ?", "user1").First(&cart).Error
	if err != nil {
		t.Errorf("Failed to retrieve cart from database: %v", err)
	}
	if len(cart.Items) != 3 {
		t.Errorf("Expected 3 items in cart after failed addition, got %d", len(cart.Items))
	}
}
