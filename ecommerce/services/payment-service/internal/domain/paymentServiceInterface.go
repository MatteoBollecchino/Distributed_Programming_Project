package domain

type PaymentServiceInterface interface {

	// Processes a payment for a given order ID and amount
	ProcessPayment(orderID string, amount float64) (*Payment, error)

	// Retrieves the payment status for a given order ID
	GetPaymentStatus(orderID string) (PaymentStatus, error)
}
