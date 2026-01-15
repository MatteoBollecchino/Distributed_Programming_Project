package domain

type CartServiceInterface interface {

	// Add an item to the cart
	AddItemToCart(username string, item CartItem) error

	// Remove an item from the cart
	RemoveItemFromCart(username string, itemID string) error

	// Update the quantity of an item in the cart
	UpdateItemQuantity(username string, itemID string, quantity int) error

	// Retrieve the cart for a user
	GetCart(username string) (*Cart, error)

	// Clear the cart for a user
	ClearCart(username string) error

	// Calculate the total price of the cart
	CalculateTotalPrice(username string) (float64, error)
}
