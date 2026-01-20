package repository

import (
	"errors"

	ulid "github.com/oklog/ulid/v2"
	"gorm.io/gorm"

	pb "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/order"
	"github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/services/order-service/internal/domain"
)

type OrderServiceRepository struct {
	db *gorm.DB
}

func NewOrderServiceRepository(db *gorm.DB) *OrderServiceRepository {
	return &OrderServiceRepository{db: db}
}

// CreateOrder creates a new order in the database.
func (r *OrderServiceRepository) CreateOrder(userID string, items []*pb.OrderItem) error {

	// Validate UserID
	if err := checkValidID(userID); err != nil {
		return err
	}

	// Validate List of Order Items
	if items == nil {
		return errors.New("items list cannot be nil")
	}
	if len(items) == 0 {
		return errors.New("order must contain at least one item")
	}
	for _, item := range items {
		if err := checkValidID(item.ItemId); err != nil {
			return err
		}
		if item.Quantity == 0 {
			return errors.New("item quantity must be greater than zero")
		}
		if item.Price < 0 {
			return errors.New("item price cannot be negative")
		}
	}

	// Check Order Uniqueness
	orderID := ulid.Make().String()
	if err := checkOrderUniqueness(r.db, orderID); err != nil {
		return err
	}

	// Create Order
	orderItems := make([]domain.OrderItem, len(items))
	for i, item := range items {
		orderItems[i] = domain.OrderItem{
			ItemID:   item.ItemId,
			Quantity: item.Quantity,
			Price:    item.Price,
		}
	}
	order := &domain.Order{
		OrderID: orderID,
		UserID:  userID,
		Items:   orderItems,
		Status:  domain.Pending,
	}

	// Save Order to Database
	if err := r.db.Create(order).Error; err != nil {
		return err
	}
	return nil
}

// UpdateOrderStatus updates the status of an existing order.
func (r *OrderServiceRepository) UpdateOrderStatus(orderID string, status pb.OrderStatus) error {

	// Validate OrderID
	if err := checkValidID(orderID); err != nil {
		return err
	}

	domainStatus, err := domain.MapProtoStatusToDomainStatus(status)
	if err != nil {
		return err
	}

	// Update Status in Database
	result := r.db.Model(&domain.Order{}).Where("order_id = ?", orderID).Update("status", domainStatus)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("order not found")
	}
	return nil
}

// GetOrder retrieves an order by its unique identifier.
func (r *OrderServiceRepository) GetOrder(orderID string) (*pb.Order, error) {

	// Validate OrderID
	if err := checkValidID(orderID); err != nil {
		return nil, err
	}

	// Retrieve Order from Database
	var domainOrder domain.Order
	if err := r.db.Preload("Items").Where("order_id = ?", orderID).First(&domainOrder).Error; err != nil {
		return nil, err
	}

	order, err := domain.DomainOrderToProtoOrder(&domainOrder)
	if err != nil {
		return nil, err
	}
	return order, nil
}

// GetOrderPrice retrieves the total price of an order by its unique identifier.
func (r *OrderServiceRepository) GetOrderPrice(orderID string) (float64, error) {

	// Validate OrderID
	if err := checkValidID(orderID); err != nil {
		return -1, err
	}

	// Retrieve Order and Calculate Total Price
	var order domain.Order
	if err := r.db.Preload("Items").Where("order_id = ?", orderID).First(&order).Error; err != nil {
		return -1, err
	}
	var totalPrice float64
	for _, item := range order.Items {
		totalPrice += float64(item.Quantity) * item.Price
	}
	return totalPrice, nil
}

// ListOrdersByUser retrieves all orders associated with a specific user.
func (r *OrderServiceRepository) ListOrdersByUser(userID string) ([]*pb.Order, error) {

	// Validate UserID
	if err := checkValidID(userID); err != nil {
		return nil, err
	}

	// Retrieve Orders from Database
	var domainOrders []*domain.Order
	if err := r.db.Preload("Items").Where("user_id = ?", userID).Find(&domainOrders).Error; err != nil {
		return nil, err
	}

	orders := make([]*pb.Order, len(domainOrders))
	for i, domainOrder := range domainOrders {
		order, err := domain.DomainOrderToProtoOrder(domainOrder)
		if err != nil {
			return nil, err
		}
		orders[i] = order
	}
	return orders, nil
}

// PRIVATE FUNCTIONS TO CHECK ON THE VALIDITY OF INPUTS

// checkValidID checks if the provided ID is valid (non-empty).
func checkValidID(id string) error {
	if id == "" {
		return errors.New("Invalid ID: cannot be empty")
	}
	return nil
}

// checkOrderUniqueness checks if the order ID is unique in the database.
// Even if orderID is generated to be unique, this function adds an extra layer of safety.
func checkOrderUniqueness(db *gorm.DB, orderID string) error {
	var count int64
	db.Model(&domain.Order{}).Where("order_id = ?", orderID).Count(&count)
	if count > 0 {
		return errors.New("Order ID already exists")
	}
	return nil
}
