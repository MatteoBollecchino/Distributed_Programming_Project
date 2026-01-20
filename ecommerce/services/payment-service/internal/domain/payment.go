package domain

import (
	"fmt"

	pb "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/payment"
)

type PaymentStatus string

const (
	PendingPayment PaymentStatus = "PENDING_PAYMENT"
	Paid           PaymentStatus = "PAID"
	PaymentFailed  PaymentStatus = "PAYMENT_FAILED"
)

type Payment struct {

	// OrderID associated with the payment
	OrderID string `gorm:"primaryKey; not null; check:order_id <> ''"`

	// Amount paid
	Amount float64 `gorm:"not null; check:amount >= 0"`

	// Current status of the payment
	Status PaymentStatus `gorm:"not null; check:status in ('PENDING_PAYMENT', 'PAID', 'PAYMENT_FAILED')"`
}

// DomainPaymentStatusToProtoPaymentStatus converts a model.Payment.Status into a pb.PaymentStatus
func DomainPaymentStatusToProtoPaymentStatus(status PaymentStatus) (pb.PaymentStatus, error) {
	switch status {
	case PendingPayment:
		return pb.PaymentStatus_PENDING_PAYMENT, nil
	case Paid:
		return pb.PaymentStatus_PAID, nil
	case PaymentFailed:
		return pb.PaymentStatus_PAYMENT_FAILED, nil
	default:
		return pb.PaymentStatus(0), fmt.Errorf("invalid domain payment status: %v", status)
	}
}
