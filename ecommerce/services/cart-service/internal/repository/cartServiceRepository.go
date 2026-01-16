package repository

import (
	"gorm.io/gorm"

	pb "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/cart"
	"github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/services/cart-service/internal/domain"
)

type CartServiceRepository struct {
	db *gorm.DB
}

func NewCartServiceRepository(db *gorm.DB) *CartServiceRepository {
	return &CartServiceRepository{db: db}
}

// Add an item to the cart
func (r *CartServiceRepository) AddItemToCart(username string, item domain.CartItem) error {
	return nil
}

// Remove an item from the cart
func (r *CartServiceRepository) RemoveItemFromCart(username string, itemID string) error {
	return nil
}

// Update the quantity of an item in the cart
func (r *CartServiceRepository) UpdateItemQuantity(username string, itemID string, quantity uint32) error {
	return nil
}

// Retrieve the cart for a user
func (r *CartServiceRepository) GetCart(username string) (*pb.Cart, error) {
	return nil, nil
}

// Clear the cart for a user
func (r *CartServiceRepository) ClearCart(username string) error {
	return nil
}

// Calculate the total price of the cart
func (r *CartServiceRepository) CalculateTotalPrice(username string) (float64, error) {
	return 0.0, nil
}
