package repository

import (
	"errors"

	"gorm.io/gorm"

	pb "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/payment"
	"github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/services/payment-service/internal/domain"
)

type PaymentServiceRepository struct {
	db *gorm.DB
}

func NewPaymentServiceRepository(db *gorm.DB) *PaymentServiceRepository {
	return &PaymentServiceRepository{db: db}
}

// CreatePayment creates a new payment for a given order ID and amount.
func (r *PaymentServiceRepository) CreatePayment(orderID string, amount float64) error {

	// Validate inputs
	if err := checkValidID(orderID); err != nil {
		return err
	}
	if amount < 0 {
		return errors.New("Invalid amount: cannot be negative")
	}

	// Check if payment already exists
	var existingPayment domain.Payment
	if err := r.db.Where("order_id = ?", orderID).First(&existingPayment).Error; err == nil {
		return errors.New("Payment already exists for this order ID")
	}

	// Create new payment with status PENDING_PAYMENT
	payment := &domain.Payment{
		OrderID: orderID,
		Amount:  amount,
		Status:  domain.PendingPayment,
	}
	if err := r.db.Create(payment).Error; err != nil {
		return err
	}
	return nil
}

// ProcessPayment processes a payment for a given order ID.
func (r *PaymentServiceRepository) ProcessPayment(orderID string, amount float64) error {
	// Validate inputs
	if err := checkValidID(orderID); err != nil {
		return err
	}

	// Retrieve the payment
	var payment domain.Payment
	if err := r.db.Where("order_id = ?", orderID).First(&payment).Error; err != nil {
		return err
	}

	// Simulate payment processing logic
	if amount >= payment.Amount {
		payment.Status = domain.Paid
	} else {
		payment.Status = domain.PaymentFailed
	}
	// Update the payment status in the database
	if err := r.db.Save(&payment).Error; err != nil {
		return err
	}
	return nil
}

// GetPaymentStatus retrieves the payment status for a given order ID.
func (r *PaymentServiceRepository) GetPaymentStatus(orderID string) (pb.PaymentStatus, error) {
	var payment domain.Payment
	if err := r.db.Where("order_id = ?", orderID).First(&payment).Error; err != nil {
		return pb.PaymentStatus(0), err
	}
	protoStatus, err := domain.DomainPaymentStatusToProtoPaymentStatus(payment.Status)
	if err != nil {
		return pb.PaymentStatus(0), err
	}
	return protoStatus, nil
}

// PRIVATE FUNCTIONS TO CHECK ON THE VALIDITY OF INPUTS

// checkValidID checks if the provided ID is valid (non-empty).
func checkValidID(id string) error {
	if id == "" {
		return errors.New("Invalid ID: cannot be empty")
	}
	return nil
}
