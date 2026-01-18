package repository

import (
	"errors"

	//"log"

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

func (r *CartServiceRepository) AddItemToCart(username string, item *pb.CartItem) error {

	// Retrieve the cart of the user from the database
	found, cart, err := r.RetrieveCart(username)
	if err != nil && found {
		return err
	}

	// If the cart does not exist, create a new one
	if !found {
		cart = &domain.Cart{
			Username: username,
			Items:    []domain.CartItem{},
		}
		if err := r.db.Session(&gorm.Session{FullSaveAssociations: true}).Create(cart).Error; err != nil {
			return err
		}
	}

	// Add the item to the cart

	// Check if the item already exists in the cart
	itemIndex := findItemInCart(cart.Items, item.ItemId)
	if itemIndex != -1 {
		cart.Items[itemIndex].Quantity += item.Quantity
	} else {
		cart.Items = append(cart.Items, domain.CartItem{
			ItemID:       item.ItemId,
			CartUsername: cart.Username,
			Quantity:     item.Quantity,
			Price:        item.Price,
		})
	}

	// Save the updated cart back to the database
	if err := r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(cart).Error; err != nil {
		return err
	}

	return nil
}

// RemoveItemFromCart removes an item from the cart of a specific user
func (r *CartServiceRepository) RemoveItemFromCart(username string, itemID string) error {

	// Retrieve the cart of the user from the database
	found, cart, err := r.RetrieveCart(username)
	if err != nil && found {
		return err
	}

	// If the cart does not exist, return an error
	if !found {
		return status.Errorf(codes.NotFound, "cart not found for user: %s", username)
	}

	// Find the item in the cart
	itemIndex := findItemInCart(cart.Items, itemID)

	// If the item is not found, return an error
	if itemIndex == -1 {
		return status.Errorf(codes.NotFound, "item with ID %s not found in cart for user: %s", itemID, username)
	}

	// Remove the item from the cart
	item := cart.Items[itemIndex]
	err = r.db.Where("cart_username = ? AND item_id = ?", item.CartUsername, item.ItemID).Delete(&domain.CartItem{}).Error
	if err != nil {
		return err
	}

	return nil
}

// UpdateItemQuantity updates the quantity of an item in the cart of a specific user
func (r *CartServiceRepository) UpdateItemQuantity(username string, itemID string, quantity uint32) error {

	// Retrieve the cart of the user from the database
	found, cart, err := r.RetrieveCart(username)
	if err != nil && found {
		return err
	}

	// If the cart does not exist, return an error
	if !found {
		return status.Errorf(codes.NotFound, "cart not found for user: %s", username)
	}

	// Find the item in the cart
	itemIndex := findItemInCart(cart.Items, itemID)

	// If the item is not found, return an error
	if itemIndex == -1 {
		return status.Errorf(codes.NotFound, "item with ID %s not found in cart for user: %s", itemID, username)
	}

	// Update the quantity of the item
	cart.Items[itemIndex].Quantity = quantity

	err = r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(cart).Error
	if err != nil {
		return err
	}

	return nil
}

// GetCart retrieves the cart for a specific user
func (r *CartServiceRepository) GetCart(username string) (*pb.Cart, error) {

	// Retrieve the cart of the user from the database
	found, cart, err := r.RetrieveCart(username)
	if err != nil && found {
		return nil, err
	}

	// If the cart does not exist, return nil and an error
	if !found {
		return nil, status.Errorf(codes.NotFound, "cart not found for user: %s", username)
	}

	// Convert the cart to pb.Cart and return it
	protoCart, err := domain.DomainCartToProtoCart(cart)
	if err != nil {
		return nil, err
	}

	return protoCart, nil
}

// ClearCart clears the cart for a specific user
func (r *CartServiceRepository) ClearCart(username string) error {

	// Retrieve the cart of the user from the database
	found, _, err := r.RetrieveCart(username)
	if err != nil && found {
		return err
	}

	// If the cart does not exist, return an error
	if !found {
		return status.Errorf(codes.NotFound, "cart not found for user: %s", username)
	}

	// Clear all items from the cart
	err = r.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&domain.CartItem{}).Error
	if err != nil {
		return err
	}

	return nil
}

// CalculateTotalPrice calculates the total price of the cart
func (r *CartServiceRepository) CalculateTotalPrice(username string) (float64, error) {

	// Retrieve the cart of the user from the database
	found, cart, err := r.RetrieveCart(username)
	if err != nil && found {
		return 0.0, err
	}

	// If the cart does not exist, return 0 and an error
	if !found {
		return 0.0, status.Errorf(codes.NotFound, "cart not found for user: %s", username)
	}

	// Calculate the total price of the items in the cart

	var total float64 = 0.0
	for _, item := range cart.Items {
		total += float64(item.Quantity) * item.Price
	}

	return total, nil
}

// RetrieveCart retrieves the cart for a specific user from the database
func (r *CartServiceRepository) RetrieveCart(username string) (bool, *domain.Cart, error) {

	var cart *domain.Cart

	err := r.db.Preload("Items").Where("username = ?", username).First(&cart).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil, status.Errorf(codes.NotFound, "cart not found")
		}
		return true, nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	return true, cart, nil
}

// findItemInCart searches for an item in the cart by its ID and returns its index and a pointer to the item
func findItemInCart(cartList []domain.CartItem, itemID string) int {

	for i, cartItem := range cartList {
		if cartItem.ItemID == itemID {
			return i
		}
	}
	return -1
}
