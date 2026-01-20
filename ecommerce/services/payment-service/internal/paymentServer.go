package internal

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/payment"
	"github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/services/payment-service/internal/domain"
)

// PaymentServer implements the payment service gRPC server.
type PaymentServer struct {
	pb.PaymentServiceServer
	repo domain.PaymentServiceInterface
}

func NewPaymentServer(repo domain.PaymentServiceInterface) *PaymentServer {
	return &PaymentServer{repo: repo}
}

// CreatePayment creates a new payment in the database.
func (s *PaymentServer) CreatePayment(ctx context.Context, req *pb.CreatePaymentRequest) (*pb.CreatePaymentResponse, error) {

	if req.OrderId == "" {
		return &pb.CreatePaymentResponse{
			ErrorMessage: "Order ID must be provided and not empty",
		}, status.Error(codes.InvalidArgument, "Order ID must be provided and not empty")
	}

	if req.Amount < 0 {
		return &pb.CreatePaymentResponse{
			ErrorMessage: "Amount cannot be negative",
		}, status.Error(codes.InvalidArgument, "Amount cannot be negative")
	}

	err := s.repo.CreatePayment(req.OrderId, req.Amount)
	if err != nil {
		return &pb.CreatePaymentResponse{ErrorMessage: err.Error()}, err
	}
	return &pb.CreatePaymentResponse{}, nil
}

// ProcessPayment processes a payment for a given order ID and amount.
func (s *PaymentServer) ProcessPayment(ctx context.Context, req *pb.ProcessPaymentRequest) (*pb.ProcessPaymentResponse, error) {

	if req.OrderId == "" {
		return &pb.ProcessPaymentResponse{
			ErrorMessage: "Order ID must be provided and not empty",
		}, status.Error(codes.InvalidArgument, "Order ID must be provided and not empty")
	}

	if req.Amount < 0 {
		return &pb.ProcessPaymentResponse{
			ErrorMessage: "Amount cannot be negative",
		}, status.Error(codes.InvalidArgument, "Amount cannot be negative")
	}

	err := s.repo.ProcessPayment(req.OrderId, req.Amount)
	if err != nil {
		return &pb.ProcessPaymentResponse{ErrorMessage: err.Error()}, err
	}
	return &pb.ProcessPaymentResponse{}, nil
}

// GetPaymentStatus retrieves the payment status for a given order ID.
func (s *PaymentServer) GetPaymentStatus(ctx context.Context, req *pb.GetPaymentStatusRequest) (*pb.GetPaymentStatusResponse, error) {

	if req.OrderId == "" {
		return &pb.GetPaymentStatusResponse{
			ErrorMessage: "Order ID must be provided and not empty",
		}, status.Error(codes.InvalidArgument, "Order ID must be provided and not empty")
	}

	status, err := s.repo.GetPaymentStatus(req.OrderId)
	if err != nil {
		return &pb.GetPaymentStatusResponse{ErrorMessage: err.Error()}, err
	}
	return &pb.GetPaymentStatusResponse{Status: status}, nil
}
