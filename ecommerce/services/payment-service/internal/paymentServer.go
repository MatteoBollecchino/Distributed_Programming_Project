package internal

/*
import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/order"
	"github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/services/order-service/internal/domain"
)

// OrderServer implements the order service gRPC server.
type OrderServer struct {
	pb.OrderServiceServer
	repo domain.OrderServiceInterface
}

func NewOrderServer(repo domain.OrderServiceInterface) *OrderServer {
	return &OrderServer{repo: repo}
}

// CreateOrder creates a new order in the database.
func (s *OrderServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {

	if req.UserId == "" {
		return &pb.CreateOrderResponse{
			ErrorMessage: "User ID must be provided and not empty",
		}, status.Error(codes.InvalidArgument, "User ID must be provided and not empty")
	}

	if req.OrderItems == nil {
		return &pb.CreateOrderResponse{
			ErrorMessage: "Items list cannot be nil",
		}, status.Error(codes.InvalidArgument, "Items list cannot be nil")
	}

	if len(req.OrderItems) == 0 {
		return &pb.CreateOrderResponse{
			ErrorMessage: "Order must contain at least one item",
		}, status.Error(codes.InvalidArgument, "Order must contain at least one item")
	}

	for _, item := range req.OrderItems {
		if item.ItemId == "" {
			return &pb.CreateOrderResponse{
				ErrorMessage: "Item ID must be provided and not empty",
			}, status.Error(codes.InvalidArgument, "Item ID must be provided and not empty")
		}
		if item.Quantity == 0 {
			return &pb.CreateOrderResponse{
				ErrorMessage: "Item quantity must be greater than zero",
			}, status.Error(codes.InvalidArgument, "Item quantity must be greater than zero")
		}
		if item.Price < 0 {
			return &pb.CreateOrderResponse{
				ErrorMessage: "Item price cannot be negative",
			}, status.Error(codes.InvalidArgument, "Item price cannot be negative")
		}
	}

	err := s.repo.CreateOrder(req.UserId, req.OrderItems)
	if err != nil {
		return &pb.CreateOrderResponse{ErrorMessage: err.Error()}, err
	}
	return &pb.CreateOrderResponse{}, nil
}

// UpdateOrderStatus updates the status of an order by its unique identifier.
func (s *OrderServer) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.UpdateOrderStatusResponse, error) {

	if req.OrderId == "" {
		return &pb.UpdateOrderStatusResponse{
			ErrorMessage: "Order ID must be provided and not empty",
		}, status.Error(codes.InvalidArgument, "Order ID must be provided and not empty")
	}

	err := s.repo.UpdateOrderStatus(req.OrderId, req.Status)
	if err != nil {
		return &pb.UpdateOrderStatusResponse{ErrorMessage: err.Error()}, err
	}
	return &pb.UpdateOrderStatusResponse{}, nil
}

// GetOrder retrieves an order by its unique identifier.
func (s *OrderServer) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {

	if req.OrderId == "" {
		return &pb.GetOrderResponse{
			ErrorMessage: "Order ID must be provided and not empty",
		}, status.Error(codes.InvalidArgument, "Order ID must be provided and not empty")
	}

	order, err := s.repo.GetOrder(req.OrderId)
	if err != nil {
		return &pb.GetOrderResponse{ErrorMessage: err.Error()}, err
	}
	return &pb.GetOrderResponse{Order: order}, nil
}

// GetOrderPrice retrieves the total price of an order by its unique identifier.
func (s *OrderServer) GetOrderPrice(ctx context.Context, req *pb.GetOrderPriceRequest) (*pb.GetOrderPriceResponse, error) {
	if req.OrderId == "" {
		return &pb.GetOrderPriceResponse{
			ErrorMessage: "Order ID must be provided and not empty",
		}, status.Error(codes.InvalidArgument, "Order ID must be provided and not empty")
	}

	totalPrice, err := s.repo.GetOrderPrice(req.OrderId)
	if err != nil {
		return &pb.GetOrderPriceResponse{ErrorMessage: err.Error()}, err
	}
	return &pb.GetOrderPriceResponse{TotalPrice: totalPrice}, nil
}

// ListOrdersByUser retrieves all orders associated with a specific user.
func (s *OrderServer) ListOrdersByUser(ctx context.Context, req *pb.ListOrdersByUserRequest) (*pb.ListOrdersByUserResponse, error) {
	if req.UserId == "" {
		return &pb.ListOrdersByUserResponse{
			ErrorMessage: "User ID must be provided and not empty",
		}, status.Error(codes.InvalidArgument, "User ID must be provided and not empty")
	}

	orders, err := s.repo.ListOrdersByUser(req.UserId)
	if err != nil {
		return &pb.ListOrdersByUserResponse{ErrorMessage: err.Error()}, err
	}
	return &pb.ListOrdersByUserResponse{Orders: orders}, nil
}*/
