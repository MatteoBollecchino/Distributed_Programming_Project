package tests

import (
	"testing"

	ulid "github.com/oklog/ulid/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	pb "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/order"
	"github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/services/order-service/internal/domain"
	"github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/services/order-service/internal/repository"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect database: %v", err)
	}

	err = db.AutoMigrate(&domain.Order{}, &domain.OrderItem{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}
	return db
}

func setupDefaultOrders(t *testing.T, db *gorm.DB) {
	defaultItem1 := &domain.OrderItem{
		ItemID:   "item123",
		Quantity: 10,
		Price:    99.99,
	}

	defaultItem2 := &domain.OrderItem{
		ItemID:   "item456",
		Quantity: 5,
		Price:    49.99,
	}

	defaultItem3 := &domain.OrderItem{
		ItemID:   "item789",
		Quantity: 2,
		Price:    19.99,
	}

	defaultItem4 := &domain.OrderItem{
		ItemID:   "item012",
		Quantity: 1,
		Price:    9.99,
	}

	order1 := &domain.Order{
		OrderID: ulid.Make().String(),
		UserID:  "user123",
		Items:   []domain.OrderItem{*defaultItem1, *defaultItem2},
		Status:  domain.Pending,
	}

	order2 := &domain.Order{
		OrderID: ulid.Make().String(),
		UserID:  "user456",
		Items:   []domain.OrderItem{*defaultItem3, *defaultItem4},
		Status:  domain.Shipped,
	}

	if err := db.Session(&gorm.Session{FullSaveAssociations: true}).Create(order1).Error; err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}
	if err := db.Session(&gorm.Session{FullSaveAssociations: true}).Create(order2).Error; err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}
}

func setupTest(t *testing.T) (*gorm.DB, *repository.OrderServiceRepository) {
	db := setupTestDB(t)
	repo := repository.NewOrderServiceRepository(db)

	setupDefaultOrders(t, db)

	return db, repo
}

func TestCreateNewOrderToNewUser(t *testing.T) {
	db, repo := setupTest(t)

	// Test creating a valid order
	err := repo.CreateOrder("user789", []*pb.OrderItem{
		{ItemId: "item111", Quantity: 3, Price: 29.99},
		{ItemId: "item222", Quantity: 1, Price: 59.99},
	})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify the order was created
	var orders []domain.Order
	if err := db.Preload("Items").Where("user_id = ?", "user789").Find(&orders).Error; err != nil {
		t.Fatalf("Failed to query orders: %v", err)
	}
	if len(orders) != 1 {
		t.Fatalf("Expected 1 order, got %d", len(orders))
	}
	if len(orders[0].Items) != 2 {
		t.Fatalf("Expected 2 items in the order, got %d", len(orders[0].Items))
	}
}

func TestCreateNewOrderToExistingUser(t *testing.T) {
	db, repo := setupTest(t)

	// Test creating a valid order for an existing user
	err := repo.CreateOrder("user123", []*pb.OrderItem{
		{ItemId: "item333", Quantity: 2, Price: 39.99},
	})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify the order was created
	var orders []domain.Order
	if err := db.Preload("Items").Where("user_id = ?", "user123").Find(&orders).Error; err != nil {
		t.Fatalf("Failed to query orders: %v", err)
	}
	if len(orders) != 2 {
		t.Fatalf("Expected 2 orders, got %d", len(orders))
	}
	if len(orders[1].Items) != 1 {
		t.Fatalf("Expected 1 item in the order, got %d", len(orders[1].Items))
	}
}

func TestCreateOrderWithInvalidUserID(t *testing.T) {
	db, repo := setupTest(t)

	// Test creating an order with empty userID
	err := repo.CreateOrder("", []*pb.OrderItem{
		{ItemId: "item444", Quantity: 1, Price: 19.99},
	})
	if err == nil {
		t.Fatalf("Expected error for empty userID, got nil")
	}

	// Verify no order was created
	var orders []domain.Order
	if err := db.Preload("Items").Where("user_id = ?", "").Find(&orders).Error; err != nil {
		t.Fatalf("Failed to query orders: %v", err)
	}
	if len(orders) != 0 {
		t.Fatalf("Expected 0 orders, got %d", len(orders))
	}
}

func TestCreateOrderWithEmptyItems(t *testing.T) {
	db, repo := setupTest(t)

	// Test creating an order with empty items
	err := repo.CreateOrder("user999", []*pb.OrderItem{})
	if err == nil {
		t.Fatalf("Expected error for empty items, got nil")
	}

	// Verify no order was created
	var orders []domain.Order
	if err := db.Preload("Items").Where("user_id = ?", "user999").Find(&orders).Error; err != nil {
		t.Fatalf("Failed to query orders: %v", err)
	}
	if len(orders) != 0 {
		t.Fatalf("Expected 0 orders, got %d", len(orders))
	}
}

func TestCreateOrderWithInvalidItemID(t *testing.T) {
	db, repo := setupTest(t)

	// Test creating an order with an invalid itemID
	err := repo.CreateOrder("user888", []*pb.OrderItem{
		{ItemId: "", Quantity: 2, Price: 29.99},
	})
	if err == nil {
		t.Fatalf("Expected error for invalid itemID, got nil")
	}

	// Verify no order was created
	var orders []domain.Order
	if err := db.Preload("Items").Where("user_id = ?", "user888").Find(&orders).Error; err != nil {
		t.Fatalf("Failed to query orders: %v", err)
	}
	if len(orders) != 0 {
		t.Fatalf("Expected 0 orders, got %d", len(orders))
	}
}

func TestCreateOrderWithInvalidQuantity(t *testing.T) {
	db, repo := setupTest(t)

	// Test creating an order with an invalid quantity
	err := repo.CreateOrder("user777", []*pb.OrderItem{
		{ItemId: "item555", Quantity: 0, Price: 39.99},
	})
	if err == nil {
		t.Fatalf("Expected error for invalid quantity, got nil")
	}

	// Verify no order was created
	var orders []domain.Order
	if err := db.Preload("Items").Where("user_id = ?", "user777").Find(&orders).Error; err != nil {
		t.Fatalf("Failed to query orders: %v", err)
	}
	if len(orders) != 0 {
		t.Fatalf("Expected 0 orders, got %d", len(orders))
	}
}

func TestCreateOrderWithInvalidPrice(t *testing.T) {
	db, repo := setupTest(t)

	// Test creating an order with an invalid price
	err := repo.CreateOrder("user666", []*pb.OrderItem{
		{ItemId: "item666", Quantity: 2, Price: -10.00},
	})
	if err == nil {
		t.Fatalf("Expected error for invalid price, got nil")
	}

	// Verify no order was created
	var orders []domain.Order
	if err := db.Preload("Items").Where("user_id = ?", "user666").Find(&orders).Error; err != nil {
		t.Fatalf("Failed to query orders: %v", err)
	}
	if len(orders) != 0 {
		t.Fatalf("Expected 0 orders, got %d", len(orders))
	}
}

func TestUpdateOrderStatusValid(t *testing.T) {
	db, repo := setupTest(t)

	// Get an existing order
	var domainOrder domain.Order
	if err := db.Preload("Items").Where("user_id = ?", "user456").First(&domainOrder).Error; err != nil {
		t.Fatalf("Failed to get existing order: %v", err)
	}

	order, err := domain.DomainOrderToProtoOrder(&domainOrder)
	if err != nil {
		t.Fatalf("Failed to convert domain order to proto order: %v", err)
	}

	// Update order status
	err = repo.UpdateOrderStatus(order.OrderId, pb.OrderStatus_SHIPPED)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify the status was updated
	var updatedOrder domain.Order
	if err := db.Where("order_id = ?", domainOrder.OrderID).First(&updatedOrder).Error; err != nil {
		t.Fatalf("Failed to get updated order: %v", err)
	}
	if updatedOrder.Status != domain.Shipped {
		t.Fatalf("Expected status %v, got %v", domain.Shipped, updatedOrder.Status)
	}
}

func TestUpdateOrderStatusInvalidID(t *testing.T) {
	db, repo := setupTest(t)

	// Attempt to update status with invalid orderID
	err := repo.UpdateOrderStatus("", pb.OrderStatus_DELIVERED)
	if err == nil {
		t.Fatalf("Expected error for invalid orderID, got nil")
	}

	// Verify no order was updated
	var orders []domain.Order
	if err := db.Preload("Items").Find(&orders).Error; err != nil {
		t.Fatalf("Failed to query orders: %v", err)
	}
	for _, order := range orders {
		if order.Status == domain.Delivered {
			t.Fatalf("Expected no orders to be updated to Delivered status")
		}
	}
}

func TestUpdateOrderStatusNonExistentID(t *testing.T) {
	db, repo := setupTest(t)

	// Attempt to update status with non-existent orderID
	err := repo.UpdateOrderStatus("nonexistentid", pb.OrderStatus_CANCELED)
	if err == nil {
		t.Fatalf("Expected error for non-existent orderID, got nil")
	}

	// Verify no order was updated
	var orders []domain.Order
	if err := db.Preload("Items").Find(&orders).Error; err != nil {
		t.Fatalf("Failed to query orders: %v", err)
	}
	for _, order := range orders {
		if order.Status == domain.Canceled {
			t.Fatalf("Expected no orders to be updated to Canceled status")
		}
	}
}

func TestGetOrderValidID(t *testing.T) {
	db, repo := setupTest(t)

	// Get an existing order
	var existingOrder domain.Order
	if err := db.Preload("Items").Where("user_id = ?", "user123").First(&existingOrder).Error; err != nil {
		t.Fatalf("Failed to get existing order: %v", err)
	}

	retrievedOrder, err := repo.GetOrder(existingOrder.OrderID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if retrievedOrder.OrderId != existingOrder.OrderID {
		t.Fatalf("Expected order ID %v, got %v", existingOrder.OrderID, retrievedOrder.OrderId)
	}
}

func TestGetOrderInvalidID(t *testing.T) {
	db, repo := setupTest(t)

	// Attempt to get order with invalid orderID
	_, err := repo.GetOrder("")
	if err == nil {
		t.Fatalf("Expected error for invalid orderID, got nil")
	}

	// Verify no order was retrieved
	var orders []domain.Order
	if err := db.Preload("Items").Find(&orders).Error; err != nil {
		t.Fatalf("Failed to query orders: %v", err)
	}
	for _, order := range orders {
		if order.OrderID == "" {
			t.Fatalf("Expected no orders to be retrieved with empty orderID")
		}
	}
}

func TestGetOrderNonExistentID(t *testing.T) {
	db, repo := setupTest(t)

	// Attempt to get order with non-existent orderID
	_, err := repo.GetOrder("nonexistentid")
	if err == nil {
		t.Fatalf("Expected error for non-existent orderID, got nil")
	}

	// Verify no order was retrieved
	var orders []domain.Order
	if err := db.Preload("Items").Find(&orders).Error; err != nil {
		t.Fatalf("Failed to query orders: %v", err)
	}
	for _, order := range orders {
		if order.OrderID == "nonexistentid" {
			t.Fatalf("Expected no orders to be retrieved with non-existent orderID")
		}
	}
}

func TestGetOrderPriceValidID(t *testing.T) {
	db, repo := setupTest(t)

	// Get an existing order
	var existingOrder domain.Order
	if err := db.Preload("Items").Where("user_id = ?", "user123").First(&existingOrder).Error; err != nil {
		t.Fatalf("Failed to get existing order: %v", err)
	}

	// Calculate expected total price
	totalPrice, err := repo.GetOrderPrice(existingOrder.OrderID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	expectedPrice := 0.0
	for _, item := range existingOrder.Items {
		expectedPrice += float64(item.Quantity) * item.Price
	}
	if totalPrice != expectedPrice {
		t.Fatalf("Expected total price %v, got %v", expectedPrice, totalPrice)
	}
}

func TestGetOrderPriceInvalidID(t *testing.T) {
	_, repo := setupTest(t)

	// Attempt to get order price with invalid orderID
	price, err := repo.GetOrderPrice("")
	if err == nil {
		t.Fatalf("Expected error for invalid orderID, got nil")
	}
	if price != -1 {
		t.Fatalf("Expected price -1 for invalid orderID, got %v", price)
	}
}

func TestGetOrderPriceNonExistentID(t *testing.T) {
	_, repo := setupTest(t)

	// Attempt to get order price with non-existent orderID
	price, err := repo.GetOrderPrice("nonexistentid")
	if err == nil {
		t.Fatalf("Expected error for non-existent orderID, got nil")
	}
	if price != -1 {
		t.Fatalf("Expected price -1 for non-existent orderID, got %v", price)
	}
}

func TestListOrdersByUserValidID(t *testing.T) {
	_, repo := setupTest(t)

	// Adding an additional order for user123 to test multiple orders
	err := repo.CreateOrder("user123", []*pb.OrderItem{
		{ItemId: "item999", Quantity: 4, Price: 14.99},
	})
	if err != nil {
		t.Fatalf("Failed to create additional order: %v", err)
	}

	// List orders for an existing user
	orders, err := repo.ListOrdersByUser("user123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(orders) != 2 {
		t.Fatalf("Expected 2 orders for user123, got %d", len(orders))
	}
}

func TestListOrdersByUserNoOrders(t *testing.T) {
	_, repo := setupTest(t)

	// List orders for a user with no orders
	orders, err := repo.ListOrdersByUser("user000")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(orders) != 0 {
		t.Fatalf("Expected 0 orders for user000, got %d", len(orders))
	}
}

func TestListOrdersByUserInvalidID(t *testing.T) {
	_, repo := setupTest(t)

	// List orders for an invalid userID (empty string)
	orders, err := repo.ListOrdersByUser("")
	if err == nil {
		t.Fatalf("Expected error for invalid userID, got nil")
	}
	if len(orders) != 0 {
		t.Fatalf("Expected 0 orders for empty userID, got %d", len(orders))
	}
}
