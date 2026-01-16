package domain

import pb "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/cart"

type CartServiceInterface interface {

	// Add an item to the cart
	AddItemToCart(username string, item *pb.CartItem) error

	// Remove an item from the cart
	RemoveItemFromCart(username string, itemID string) error

	// Update the quantity of an item in the cart
	UpdateItemQuantity(username string, itemID string, quantity uint32) error

	// Retrieve the cart for a user
	GetCart(username string) (*pb.Cart, error)

	// Clear the cart for a user
	ClearCart(username string) error

	// Calculate the total price of the cart
	CalculateTotalPrice(username string) (float64, error)
}
