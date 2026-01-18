package tests

import (
	"log"
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

func setupDefaultCarts(t *testing.T, db *gorm.DB) {

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

	setupDefaultCarts(t, db)

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

func TestRemoveExistingItemFromCart(t *testing.T) {
	db, repo := setupTest(t)

	// Verify initial quantity of item3 in user2's cart
	var initialCart domain.Cart
	if err := db.Preload("Items").Where("username = ?", "user2").First(&initialCart).Error; err != nil {
		t.Fatalf("Failed to retrieve initial cart: %v", err)
	}

	if len(initialCart.Items) != 2 {
		t.Fatalf("Expected 2 items before removal, got %d", len(initialCart.Items))
	}

	// Test removing an existing item from the cart
	err := repo.RemoveItemFromCart("user2", "item3")
	if err != nil {
		t.Fatalf("Failed to remove existing item from cart: %v", err)
	}

	// Verify if the item was removed correctly
	var cart domain.Cart
	if err := db.Preload("Items").Where("username = ?", "user2").First(&cart).Error; err != nil {
		t.Fatalf("Failed to retrieve cart from database: %v", err)
	}

	if len(cart.Items) != 1 {
		t.Fatalf("Expected 1 item in cart after removal, got %d", len(cart.Items))
	}

	log.Printf("%v", cart.Items)

	for _, item := range cart.Items {
		if item.ItemID == "item3" {
			t.Fatalf("item3 should have been removed from cart")
		}
	}
}

func TestRemoveNonExistingItemFromCart(t *testing.T) {
	db, repo := setupTest(t)

	// Test removing a non-existing item from the cart
	err := repo.RemoveItemFromCart("user1", "nonexistent_item")
	if err == nil {
		t.Errorf("Expected error when removing non-existing item, got nil")
	}

	// Verify that the cart remains unchanged
	var cart domain.Cart
	err = db.Preload("Items").Where("username = ?", "user1").First(&cart).Error
	if err != nil {
		t.Errorf("Failed to retrieve cart from database: %v", err)
	}
	if len(cart.Items) != 2 {
		t.Errorf("Expected 2 items in cart after failed removal, got %d", len(cart.Items))
	}
}

func TestRemoveItemFromNonExistingCart(t *testing.T) {
	db, repo := setupTest(t)

	// Test removing an item from a non-existing cart
	err := repo.RemoveItemFromCart("nonexistent_user", "item1")
	if err == nil {
		t.Errorf("Expected error when removing item from non-existing cart, got nil")
	}

	// Verify that no cart was created
	var cart domain.Cart
	err = db.Preload("Items").Where("username = ?", "nonexistent_user").First(&cart).Error
	if err == nil {
		t.Errorf("Expected no cart for nonexistent_user, but found one")
	}
}

func TestRemoveItemFromCartWithEmptyID(t *testing.T) {
	db, repo := setupTest(t)

	// Test removing an item with empty ID
	err := repo.RemoveItemFromCart("user1", "")
	if err == nil {
		t.Errorf("Expected error when removing item with empty ID, got nil")
	}

	// Verify that the cart remains unchanged
	var cart domain.Cart
	err = db.Preload("Items").Where("username = ?", "user1").First(&cart).Error
	if err != nil {
		t.Errorf("Failed to retrieve cart from database: %v", err)
	}
	if len(cart.Items) != 2 {
		t.Errorf("Expected 2 items in cart after failed removal, got %d", len(cart.Items))
	}
}

func TestRemoveItemFromCartWithEmptyUsername(t *testing.T) {
	db, repo := setupTest(t)

	// Test removing an item with empty username
	err := repo.RemoveItemFromCart("", "item1")
	if err == nil {
		t.Errorf("Expected error when removing item with empty username, got nil")
	}

	// Verify that the cart remains unchanged
	var cart domain.Cart
	err = db.Preload("Items").Where("username = ?", "").First(&cart).Error
	if err == nil {
		t.Errorf("Failed to retrieve cart from database: %v", err)
	}
	if len(cart.Items) != 0 {
		t.Errorf("Expected 0 items in cart after failed removal, got %d", len(cart.Items))
	}
}

func TestRemoveItemFromCartWithEmptyIDAndUsername(t *testing.T) {
	db, repo := setupTest(t)

	// Test removing an item with empty ID and username
	err := repo.RemoveItemFromCart("", "")
	if err == nil {
		t.Errorf("Expected error when removing item with empty ID and username, got nil")
	}

	// Verify that the cart remains unchanged
	var cart domain.Cart
	err = db.Preload("Items").Where("username = ?", "").First(&cart).Error
	if err == nil {
		t.Errorf("Failed to retrieve cart from database: %v", err)
	}
	if len(cart.Items) != 0 {
		t.Errorf("Expected 0 items in cart after failed removal, got %d", len(cart.Items))
	}
}

func TestUpdateItemQuantityInCart(t *testing.T) {
	db, repo := setupTest(t)

	// Test updating the quantity of an existing item in the cart
	err := repo.UpdateItemQuantity("user1", "item1", 5)
	if err != nil {
		t.Errorf("Failed to update item quantity in cart: %v", err)
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

func TestUpdateNonExistingItemQuantityInCart(t *testing.T) {
	db, repo := setupTest(t)

	var previousCart domain.Cart
	err := db.Preload("Items").Where("username = ?", "user1").First(&previousCart).Error

	// Test updating the quantity of a non-existing item in the cart
	err = repo.UpdateItemQuantity("user1", "nonexistent_item", 5)
	if err == nil {
		t.Errorf("Expected error when updating quantity of non-existing item, got nil")
	}

	// Verify that the cart remains unchanged
	var cart domain.Cart
	err = db.Preload("Items").Where("username = ?", "user1").First(&cart).Error
	if err != nil {
		t.Errorf("Failed to retrieve cart from database: %v", err)
	}

	for i, item := range cart.Items {
		if item.Quantity != previousCart.Items[i].Quantity {
			t.Errorf("Expected quantity %d for item%d, got %d", previousCart.Items[i].Quantity, i+1, item.Quantity)
		}
	}
}

func TestUpdateItemQuantityInNonExistingCart(t *testing.T) {
	db, repo := setupTest(t)

	// Test updating the quantity of an item in a non-existing cart
	err := repo.UpdateItemQuantity("nonexistent_user", "item1", 5)
	if err == nil {
		t.Errorf("Expected error when updating item quantity in non-existing cart, got nil")
	}

	// Verify that no cart was created
	var cart domain.Cart
	err = db.Preload("Items").Where("username = ?", "nonexistent_user").First(&cart).Error
	if err == nil {
		t.Errorf("Expected no cart for nonexistent_user, but found one")
	}
}

func TestUpdateItemQuantityInCartWithEmptyID(t *testing.T) {
	db, repo := setupTest(t)

	var previousCart domain.Cart
	err := db.Preload("Items").Where("username = ?", "user1").First(&previousCart).Error

	// Test updating the quantity of an item with empty ID
	err = repo.UpdateItemQuantity("user1", "", 5)
	if err == nil {
		t.Errorf("Expected error when updating item quantity with empty ID, got nil")
	}

	// Verify that the cart remains unchanged
	var cart domain.Cart
	err = db.Preload("Items").Where("username = ?", "user1").First(&cart).Error
	if err != nil {
		t.Errorf("Failed to retrieve cart from database: %v", err)
	}

	for i, item := range cart.Items {
		if item.Quantity != previousCart.Items[i].Quantity {
			t.Errorf("Expected quantity %d for item%d, got %d", previousCart.Items[i].Quantity, i+1, item.Quantity)
		}
	}
}

func TestUpdateItemQuantityInCartWithEmptyUsername(t *testing.T) {
	db, repo := setupTest(t)

	// Test updating the quantity of an item with empty username
	err := repo.UpdateItemQuantity("", "item1", 5)
	if err == nil {
		t.Errorf("Expected error when updating item quantity with empty username, got nil")
	}

	// Verify that the cart remains unchanged
	var cart domain.Cart
	err = db.Preload("Items").Where("username = ?", "").First(&cart).Error
	if err == nil {
		t.Errorf("Failed to retrieve cart from database: %v", err)
	}
	if len(cart.Items) != 0 {
		t.Errorf("Expected 0 items in cart after failed update, got %d", len(cart.Items))
	}
}

func TestUpdateItemQuantityInCartWithEmptyIDAndUsername(t *testing.T) {
	db, repo := setupTest(t)

	// Test updating the quantity of an item with empty ID and username
	err := repo.UpdateItemQuantity("", "", 5)
	if err == nil {
		t.Errorf("Expected error when updating item quantity with empty ID and username, got nil")
	}

	// Verify that the cart remains unchanged
	var cart domain.Cart
	err = db.Preload("Items").Where("username = ?", "").First(&cart).Error
	if err == nil {
		t.Errorf("Failed to retrieve cart from database: %v", err)
	}
	if len(cart.Items) != 0 {
		t.Errorf("Expected 0 items in cart after failed update, got %d", len(cart.Items))
	}
}

func TestGetExistingCart(t *testing.T) {
	_, repo := setupTest(t)

	// Test retrieving an existing cart
	cart, err := repo.GetCart("user1")
	if err != nil {
		t.Errorf("Failed to retrieve existing cart: %v", err)
	}
	if cart.Username != "user1" {
		t.Errorf("Expected cart username 'user1', got '%s'", cart.Username)
	}
	if len(cart.Items) != 2 {
		t.Errorf("Expected 2 items in cart, got %d", len(cart.Items))
	}
}

func TestGetNonExistingCart(t *testing.T) {
	_, repo := setupTest(t)

	// Test retrieving a non-existing cart
	_, err := repo.GetCart("nonexistent_user")
	if err == nil {
		t.Errorf("Expected error when retrieving non-existing cart, got nil")
	}
}

func TestGetCartWithEmptyUsername(t *testing.T) {
	_, repo := setupTest(t)

	// Test retrieving a cart with empty username
	_, err := repo.GetCart("")
	if err == nil {
		t.Errorf("Expected error when retrieving cart with empty username, got nil")
	}
}

func TestClearExistingCart(t *testing.T) {
	db, repo := setupTest(t)

	// Test clearing an existing cart
	err := repo.ClearCart("user2")
	if err != nil {
		t.Errorf("Failed to clear existing cart: %v", err)
	}

	// Verify if the cart was cleared correctly
	var cart domain.Cart
	err = db.Preload("Items").Where("username = ?", "user2").First(&cart).Error
	if err != nil {
		t.Errorf("Failed to retrieve cart from database: %v", err)
	}
	if len(cart.Items) != 0 {
		t.Errorf("Expected 0 items in cart after clearing, got %d", len(cart.Items))
	}
}

func TestClearNonExistingCart(t *testing.T) {
	_, repo := setupTest(t)

	// Test clearing a non-existing cart
	err := repo.ClearCart("nonexistent_user")
	if err == nil {
		t.Errorf("Expected error when clearing non-existing cart, got nil")
	}
}

func TestClearCartWithEmptyUsername(t *testing.T) {
	_, repo := setupTest(t)

	// Test clearing a cart with empty username
	err := repo.ClearCart("")
	if err == nil {
		t.Errorf("Expected error when clearing cart with empty username, got nil")
	}
}

// DA CONTROLLARE DA QUI IN GIU'

func TestCalculateTotalPriceOfExistingCart(t *testing.T) {
	_, repo := setupTest(t)
	// Test calculating the total price of an existing cart
	total, err := repo.CalculateTotalPrice("user1")
	if err != nil {
		t.Errorf("Failed to calculate total price of existing cart: %v", err)
	}
	expectedTotal := 2*10.0 + 1*20.0 // item1 + item2
	if total != expectedTotal {
		t.Errorf("Expected total price %.2f, got %.2f", expectedTotal, total)
	}
}

func TestCalculateTotalPriceOfNonExistingCart(t *testing.T) {
	_, repo := setupTest(t)

	// Test calculating the total price of a non-existing cart
	_, err := repo.CalculateTotalPrice("nonexistent_user")
	if err == nil {
		t.Errorf("Expected error when calculating total price of non-existing cart, got nil")
	}
}

func TestCalculateTotalPriceOfCartWithEmptyUsername(t *testing.T) {
	_, repo := setupTest(t)
	// Test calculating the total price of a cart with empty username
	_, err := repo.CalculateTotalPrice("")
	if err == nil {
		t.Errorf("Expected error when calculating total price of cart with empty username, got nil")
	}
}
