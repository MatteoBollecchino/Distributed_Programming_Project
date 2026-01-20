package tests

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/services/payment-service/internal/domain"
	"github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/services/payment-service/internal/repository"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect database: %v", err)
	}

	err = db.AutoMigrate(&domain.Payment{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}
	return db
}

func setupDefaultPayments(t *testing.T, db *gorm.DB) {
	defaultPayment1 := &domain.Payment{
		OrderID: "order123",
		Amount:  199.99,
		Status:  domain.PendingPayment,
	}

	defaultPayment2 := &domain.Payment{
		OrderID: "order456",
		Amount:  49.99,
		Status:  domain.Paid,
	}

	defaultPayment3 := &domain.Payment{
		OrderID: "order789",
		Amount:  39.99,
		Status:  domain.PaymentFailed,
	}

	if err := db.Create(defaultPayment1).Error; err != nil {
		t.Fatalf("Failed to create payment: %v", err)
	}
	if err := db.Create(defaultPayment2).Error; err != nil {
		t.Fatalf("Failed to create payment: %v", err)
	}
	if err := db.Create(defaultPayment3).Error; err != nil {
		t.Fatalf("Failed to create payment: %v", err)
	}
}

func setupTest(t *testing.T) (*gorm.DB, *repository.PaymentServiceRepository) {
	db := setupTestDB(t)
	repo := repository.NewPaymentServiceRepository(db)

	setupDefaultPayments(t, db)

	return db, repo
}
