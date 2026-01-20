package domain

import (
	"fmt"

	pb "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/order"
)

type OrderItem struct {

	// OrderID is the unique identifier for the order to which the item belongs.
	OrderID string `gorm:"not null; check:order_id <> ''"`

	// ItemID is the unique identifier for the item.
	ItemID string `gorm:"not null; check:item_id <> ''"`

	// Quantity indicates the number of units of the item in the order.
	Quantity uint32 `gorm:"not null; check:quantity > 0"`

	// Price represents the price of a single unit of the item at the time of the order.
	Price float64 `gorm:"not null; check:price >= 0"`
}

// DomainOrderItemToProtoOrderItem converts a model.OrderItem into a pb.OrderItem
func DomainOrderItemToProtoOrderItem(item *OrderItem) (*pb.OrderItem, error) {
	if item == nil {
		return nil, fmt.Errorf("Input argument is nil")
	}
	return &pb.OrderItem{
		ItemId:   item.ItemID,
		Quantity: item.Quantity,
		Price:    item.Price,
	}, nil
}
