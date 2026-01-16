package domain

import (
	"fmt"

	pb "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/cart"
)

type CartItem struct {

	// ItemID is the unique identifier for the item.
	ItemID string `gorm:"primaryKey"`

	// Quantity indicates how many of the item are in the cart.
	Quantity uint32

	// Price indicates the price of a single item.
	Price float64
}

// DomainCartItemToProtoCartItem converts a model.CartItem into a pb.CartItem
func DomainCartItemToProtoCartItem(cartItem *CartItem) (*pb.CartItem, error) {
	if cartItem == nil {
		return nil, fmt.Errorf("Input argument is nil")
	}

	return &pb.CartItem{
		ItemId:   cartItem.ItemID,
		Quantity: uint32(cartItem.Quantity),
		Price:    cartItem.Price,
	}, nil
}

/*
// ProtoCartItemToDomainCartItem converts a pb.CartItem into a model.CartItem
func ProtoCartItemToDomainCartItem(protoCartItem *pb.CartItem) (*CartItem, error) {
	if protoCartItem == nil {
		return nil, fmt.Errorf("Input argument is nil")
	}

	return &CartItem{
		ItemID:   protoCartItem.ItemId,
		Quantity: uint32(protoCartItem.Quantity),
		Price:    protoCartItem.Price,
	}, nil
}*/
