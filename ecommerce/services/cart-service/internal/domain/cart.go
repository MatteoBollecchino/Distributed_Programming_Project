package domain

type Cart struct {

	// Username is the unique identifier for the cart owner.
	Username string

	// Items holds the items in the cart.
	Items []CartItem
}

type CartItem struct {

	// ItemID is the unique identifier for the item.
	ItemID string

	// Quantity indicates how many of the item are in the cart.
	Quantity int
}

/*
// DomainCartToProtoCart converts a model.Cart into a pb.Cart
func DomainCartToProtoCart(cart *Cart) (*pb.Cart, error) {
	if cart == nil {
		return nil, fmt.Errorf("Input argument is nil")
	}

	var r pb.Role
	if user.Role == AdminRole {
		r = pb.Role_ADMIN
	} else {
		r = pb.Role_USER
	}
	return &pb.User{
		Username: user.Username,
		Password: user.Password,
		Role:     r}, nil
}*/
