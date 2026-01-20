package tests

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	pb "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/payment"
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

func TestCreateNewPayment(t *testing.T) {
	db, repo := setupTest(t)

	err := repo.CreatePayment("order999", 59.99)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var payment domain.Payment
	if err := db.Where("order_id = ?", "order999").First(&payment).Error; err != nil {
		t.Fatalf("Failed to retrieve payment: %v", err)
	}
	if payment.Amount != 59.99 || payment.Status != domain.PendingPayment {
		t.Fatalf("Payment data mismatch: got %+v", payment)
	}
}

func TestCreatePaymentAlreadyExists(t *testing.T) {
	db, repo := setupTest(t)

	err := repo.CreatePayment("order123", 199.99)
	if err == nil {
		t.Fatalf("Expected error for existing payment, got nil")
	}

	var count int64
	if err := db.Model(&domain.Payment{}).Where("order_id = ?", "order123").Count(&count).Error; err != nil {
		t.Fatalf("Failed to count payments: %v", err)
	}
	if count != 1 {
		t.Fatalf("Expected 1 payment record, got %d", count)
	}
}

func TestCreatePaymentInvalidAmount(t *testing.T) {
	db, repo := setupTest(t)

	err := repo.CreatePayment("order456", -10.00)
	if err == nil {
		t.Fatalf("Expected error for negative amount, got nil")
	}

	var count int64
	if err := db.Model(&domain.Payment{}).Where("order_id = ?", "order456").Count(&count).Error; err != nil {
		t.Fatalf("Failed to count payments: %v", err)
	}
	if count != 1 {
		t.Fatalf("Expected 1 payment record, got %d", count)
	}
}

func TestProcessPayment(t *testing.T) {
	db, repo := setupTest(t)

	err := repo.ProcessPayment("order123", 199.99)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var payment domain.Payment
	if err := db.Where("order_id = ?", "order123").First(&payment).Error; err != nil {
		t.Fatalf("Failed to retrieve payment: %v", err)
	}
	if payment.Status != domain.Paid {
		t.Fatalf("Expected payment status to be PAID, got %v", payment.Status)
	}
}

func TestProcessPaymentInvalidID(t *testing.T) {
	_, repo := setupTest(t)

	err := repo.ProcessPayment("", 50.00)
	if err == nil {
		t.Fatalf("Expected error for invalid order ID, got nil")
	}
}

func TestProcessPaymentNegativeAmount(t *testing.T) {
	_, repo := setupTest(t)

	err := repo.ProcessPayment("order123", -20.00)
	if err == nil {
		t.Fatalf("Expected error for negative amount, got nil")
	}
}

func TestProcessPaymentInsufficientAmount(t *testing.T) {
	db, repo := setupTest(t)

	err := repo.ProcessPayment("order123", 100.00)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var payment domain.Payment
	if err := db.Where("order_id = ?", "order123").First(&payment).Error; err != nil {
		t.Fatalf("Failed to retrieve payment: %v", err)
	}
	if payment.Status != domain.PaymentFailed {
		t.Fatalf("Expected payment status to be PAYMENT_FAILED, got %v", payment.Status)
	}
}

func TestProcessPaymentAlreadyPaid(t *testing.T) {
	db, repo := setupTest(t)

	err := repo.ProcessPayment("order456", 49.99)
	if err == nil {
		t.Fatalf("Expected error for already PAID payment, got nil")
	}

	var payment domain.Payment
	if err := db.Where("order_id = ?", "order456").First(&payment).Error; err != nil {
		t.Fatalf("Failed to retrieve payment: %v", err)
	}
	if payment.Status != domain.Paid {
		t.Fatalf("Expected payment status to remain PAID, got %v", payment.Status)
	}
}

func TestProcessPaymentFailed(t *testing.T) {
	db, repo := setupTest(t)

	err := repo.ProcessPayment("order789", 40.00)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var payment domain.Payment
	if err := db.Where("order_id = ?", "order789").First(&payment).Error; err != nil {
		t.Fatalf("Failed to retrieve payment: %v", err)
	}
	if payment.Status == domain.PaymentFailed {
		t.Fatalf("Expected payment status to be PAYMENT_FAILED, got %v", payment.Status)
	}
	if payment.Status != domain.Paid {
		t.Fatalf("Expected payment status to be PAID, got %v", payment.Status)
	}
}

func TestGetPaymentStatus(t *testing.T) {
	_, repo := setupTest(t)

	status, err := repo.GetPaymentStatus("order456")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if status != pb.PaymentStatus_PAID {
		t.Fatalf("Expected payment status to be PAID, got %v", status)
	}

	status, err = repo.GetPaymentStatus("order789")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if status != pb.PaymentStatus_PAYMENT_FAILED {
		t.Fatalf("Expected payment status to be PAYMENT_FAILED, got %v", status)
	}

	status, err = repo.GetPaymentStatus("order123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if status != pb.PaymentStatus_PENDING_PAYMENT {
		t.Fatalf("Expected payment status to be PENDING_PAYMENT, got %v", status)
	}
}

func TestGetPaymentStatusNonExistentOrder(t *testing.T) {
	_, repo := setupTest(t)

	_, err := repo.GetPaymentStatus("nonexistent_order")
	if err == nil {
		t.Fatalf("Expected error for nonexistent order ID, got nil")
	}
}

func TestGetPaymentStatusInvalidIDFormat(t *testing.T) {
	_, repo := setupTest(t)

	_, err := repo.GetPaymentStatus("")
	if err == nil {
		t.Fatalf("Expected error for invalid order ID format, got nil")
	}
}
