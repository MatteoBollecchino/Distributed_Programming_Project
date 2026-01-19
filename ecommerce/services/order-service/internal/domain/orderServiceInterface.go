package domain

import pb "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/order"

type OrderServiceInterface interface {

	// CreateOrder creates a new order with the provided details.
	CreateOrder(userID string, items []pb.OrderItem) error

	// UpdateOrderStatus updates the status of an order by its unique identifier.
	UpdateOrderStatus(orderID string, status Status) error

	// GetOrder retrieves an order by its unique identifier.
	GetOrder(orderID string) (*pb.Order, error)

	// GetOrderPrice retrieves the total price of an order by its unique identifier.
	GetOrderPrice(orderID string) (float64, error)

	// ListOrdersByUser retrieves all orders associated with a specific user.
	ListOrdersByUser(userID string) ([]*pb.Order, error)
}
