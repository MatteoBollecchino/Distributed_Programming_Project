package domain

import (
	"fmt"

	pb "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/order"
)

type Status string

const (
	// Pending indicates that the order has been created but not yet processed.
	Pending Status = "PENDING"

	// Processing indicates that the order is currently being processed.
	Processing Status = "PROCESSING"

	// Shipped indicates that the order has been shipped to the customer.
	Shipped Status = "SHIPPED"

	// Delivered indicates that the order has been delivered to the customer.
	Delivered Status = "DELIVERED"

	// Canceled indicates that the order has been canceled.
	Canceled Status = "CANCELED"
)

type Order struct {

	// OrderID is the unique identifier for the order.
	OrderID string `gorm:"primaryKey; not null; check:order_id <> ''"`

	// UserID is the unique identifier for the user who placed the order.
	UserID string `gorm:"not null; check:user_id <> ''"`

	// ItemIDs is a list of unique identifiers for the items included in the order.
	Items []OrderItem `gorm:"foreignKey:ItemID;references:OrderID;constraint:OnDelete:CASCADE;not null"`

	// Status indicates the current status of the order (e.g., "Pending", "Shipped", "Delivered").
	Status Status `gorm:"not null; check:status <> ''"`
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

// DomainOrderToProtoOrder converts a model.Order into a pb.Order
func DomainOrderToProtoOrder(order *Order) (*pb.Order, error) {
	if order == nil {
		return nil, fmt.Errorf("Input argument is nil")
	}

	return &pb.Order{
		OrderId: order.OrderID,
		UserId:  order.UserID,
		Items:   []*pb.OrderItem{},
		Status:  pb.OrderStatus(pb.OrderStatus_value[string(order.Status)]),
	}, nil
}
