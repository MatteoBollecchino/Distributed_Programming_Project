package tests

import (
	"log"
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
	log.Println(orders[0].Items[0].ItemID)
	if len(orders) != 2 {
		t.Fatalf("Expected 2 orders, got %d", len(orders))
	}
	if len(orders[1].Items) != 1 {
		t.Fatalf("Expected 1 item in the order, got %d", len(orders[1].Items))
	}
}
