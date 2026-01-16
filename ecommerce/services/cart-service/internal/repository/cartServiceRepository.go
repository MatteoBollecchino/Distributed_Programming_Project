package repository

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

// AddItemToCart adds an item to the cart of a specific user
func (r *CartServiceRepository) AddItemToCart(username string, item *pb.CartItem) error {

	// Retrieve the cart of the user from the database
	found, cart, err := r.RetrieveCart(username)
	if err != nil {
		return err
	}

	// If the cart does not exist, create a new one
	if !found {
		err := r.db.Create(&cart).Error
		if err != nil {

		}
	}

	// Add the item to the cart

	// Save the updated cart back to the database

	return nil
}

// RemoveItemFromCart removes an item from the cart of a specific user
func (r *CartServiceRepository) RemoveItemFromCart(username string, itemID string) error {

	// Retrieve the cart of the user from the database

	// If the cart does not exist, return an error

	// Find the item in the cart

	// Remove the item from the cart
	return nil
}

// UpdateItemQuantity updates the quantity of an item in the cart of a specific user
func (r *CartServiceRepository) UpdateItemQuantity(username string, itemID string, quantity uint32) error {

	// Retrieve the cart of the user from the database

	// If the cart does not exist, return an error

	// Find the item in the cart

	// Update the quantity of the item

	return nil
}

// GetCart retrieves the cart for a specific user
func (r *CartServiceRepository) GetCart(username string) (*pb.Cart, error) {

	// Retrieve the cart of the user from the database

	// If the cart does not exist, return nil and an error

	// Convert the cart to pb.Cart and return it

	return nil, nil
}

// ClearCart clears the cart for a specific user
func (r *CartServiceRepository) ClearCart(username string) error {

	// Retrieve the cart of the user from the database

	// If the cart does not exist, return an error

	// Clear all items from the cart

	return nil
}

// CalculateTotalPrice calculates the total price of the cart
func (r *CartServiceRepository) CalculateTotalPrice(username string) (float64, error) {

	// Retrieve the cart of the user from the database

	// If the cart does not exist, return 0 and an error

	// Calculate the total price of the items in the cart

	return 0.0, nil
}

// RetrieveCart retrieves the cart for a specific user from the database
func (r *CartServiceRepository) RetrieveCart(username string) (bool, *domain.Cart, error) {

	var cart *domain.Cart

	err := r.db.Where("username = ?", username).First(&cart).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil, status.Errorf(codes.NotFound, "cart not found")
		}
		return true, nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	return true, cart, nil
}
