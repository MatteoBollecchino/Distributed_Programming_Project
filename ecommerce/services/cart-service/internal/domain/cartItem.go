package domain

import (
	"fmt"

	pb "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/cart"
)

type CartItem struct {

	// ItemID is the unique identifier for the item.
	ItemID string

	// Quantity indicates how many of the item are in the cart.
	Quantity uint32
}

// DomainCartItemToProtoCartItem converts a model.CartItem into a pb.CartItem
func DomainCartItemToProtoCartItem(cartItem *CartItem) (*pb.CartItem, error) {
	if cartItem == nil {
		return nil, fmt.Errorf("Input argument is nil")
	}

	return &pb.CartItem{
		ItemId:   cartItem.ItemID,
		Quantity: uint32(cartItem.Quantity),
	}, nil
}
