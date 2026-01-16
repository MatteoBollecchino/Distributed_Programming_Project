package domain

import (
	"fmt"

	pb "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/cart"
)

type Cart struct {

	// Username is the unique identifier for the cart owner.
	Username string `gorm:"primaryKey"`

	// Items holds the items in the cart.
	Items []CartItem `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE"`
}

// DomainCartToProtoCart converts a model.Cart into a pb.Cart
func DomainCartToProtoCart(cart *Cart) (*pb.Cart, error) {
	if cart == nil {
		return nil, fmt.Errorf("Input argument is nil")
	}

	protoCart := &pb.Cart{
		Username: cart.Username,
		Items:    []*pb.CartItem{},
	}
	for _, item := range cart.Items {
		protoItem, err := DomainCartItemToProtoCartItem(&item)
		if err != nil {
			return nil, err
		}
		protoCart.Items = append(protoCart.Items, protoItem)
	}
	return protoCart, nil
}
