package domain

type PaymentStatus string

const (
	PendingPayment PaymentStatus = "PENDING_PAYMENT"
	Paid           PaymentStatus = "PAID"
	PaymentFailed  PaymentStatus = "PAYMENT_FAILED"
)

type Payment struct {
	OrderID string
	Amount  float64
	Status  PaymentStatus
}
