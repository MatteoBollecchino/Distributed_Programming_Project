package domain

import pb "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/payment"

type PaymentServiceInterface interface {

	// Creates a new payment
	CreatePayment(orderID string, amount float64) error

	// Processes a payment for a given order ID and amount
	ProcessPayment(orderID string, amount float64) error

	// Retrieves the payment status for a given order ID
	GetPaymentStatus(orderID string) (pb.PaymentStatus, error)
}
