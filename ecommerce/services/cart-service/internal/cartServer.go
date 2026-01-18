package internal

import (
	"context"

	pb "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/cart"
	"github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/services/cart-service/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CartServer implements the cart service gRPC server.
type CartServer struct {
	pb.CartServiceServer
	repo domain.CartServiceInterface
}

func NewCartServer(repo domain.CartServiceInterface) *CartServer {
	return &CartServer{repo: repo}
}

// AddItemToCart adds an item to the cart of a specific user.
func (s *CartServer) AddItemToCart(ctx context.Context, req *pb.AddItemToCartRequest) (*pb.AddItemToCartResponse, error) {

	if req.Username == "" || req.CartItem == nil {
		return &pb.AddItemToCartResponse{
			ErrorMessage: "Username and CartItem must be provided and not empty or nil",
		}, status.Error(codes.InvalidArgument, "Username and CartItem must be provided and not empty or nil")
	}

	if req.CartItem.Quantity == 0 {
		return &pb.AddItemToCartResponse{
			ErrorMessage: "Quantity must be greater than zero",
		}, status.Error(codes.InvalidArgument, "Quantity must be greater than zero")
	}

	err := s.repo.AddItemToCart(req.Username, req.CartItem)
	if err != nil {
		return &pb.AddItemToCartResponse{ErrorMessage: err.Error()}, err
	}
	return &pb.AddItemToCartResponse{}, nil
}

// RemoveItemFromCart removes an item from the cart of a specific user
func (s *CartServer) RemoveItemFromCart(ctx context.Context, req *pb.RemoveItemFromCartRequest) (*pb.RemoveItemFromCartResponse, error) {

	if req.Username == "" || req.ItemId == "" {
		return &pb.RemoveItemFromCartResponse{
			ErrorMessage: "Username and ItemId must be provided and not empty",
		}, status.Error(codes.InvalidArgument, "Username and ItemId must be provided and not empty")
	}

	err := s.repo.RemoveItemFromCart(req.Username, req.ItemId)
	if err != nil {
		return &pb.RemoveItemFromCartResponse{ErrorMessage: err.Error()}, err
	}
	return &pb.RemoveItemFromCartResponse{}, nil
}

// UpdateItemQuantity updates the quantity of an item in the cart of a specific user
func (s *CartServer) UpdateItemQuantity(ctx context.Context, req *pb.UpdateItemQuantityRequest) (*pb.UpdateItemQuantityResponse, error) {

	if req.Username == "" || req.ItemId == "" {
		return &pb.UpdateItemQuantityResponse{
			ErrorMessage: "Username and ItemId must be provided and not empty",
		}, status.Error(codes.InvalidArgument, "Username and ItemId must be provided and not empty")
	}

	if req.Quantity == 0 {
		return &pb.UpdateItemQuantityResponse{
			ErrorMessage: "Quantity must be greater than zero",
		}, status.Error(codes.InvalidArgument, "Quantity must be greater than zero")
	}

	err := s.repo.UpdateItemQuantity(req.Username, req.ItemId, req.Quantity)
	if err != nil {
		return &pb.UpdateItemQuantityResponse{ErrorMessage: err.Error()}, err
	}
	return &pb.UpdateItemQuantityResponse{}, nil
}

// GetCart retrieves the cart for a specific user
func (s *CartServer) GetCart(ctx context.Context, req *pb.GetCartRequest) (*pb.GetCartResponse, error) {

	if req.Username == "" {
		return &pb.GetCartResponse{
			Cart:         nil,
			ErrorMessage: "Username must be provided and not empty",
		}, status.Error(codes.InvalidArgument, "Username must be provided and not empty")
	}

	cart, err := s.repo.GetCart(req.Username)
	if err != nil {
		return &pb.GetCartResponse{Cart: nil, ErrorMessage: err.Error()}, err
	}
	return &pb.GetCartResponse{Cart: cart}, nil
}

// ClearCart clears the cart for a specific user
func (s *CartServer) ClearCart(ctx context.Context, req *pb.ClearCartRequest) (*pb.ClearCartResponse, error) {

	if req.Username == "" {
		return &pb.ClearCartResponse{
			ErrorMessage: "Username must be provided and not empty",
		}, status.Error(codes.InvalidArgument, "Username must be provided and not empty")
	}

	err := s.repo.ClearCart(req.Username)
	if err != nil {
		return &pb.ClearCartResponse{ErrorMessage: err.Error()}, err
	}

	return &pb.ClearCartResponse{}, nil
}

// CalculateTotalPrice calculates the total price of the cart
func (s *CartServer) CalculateTotalPrice(ctx context.Context, req *pb.CalculateTotalPriceRequest) (*pb.CalculateTotalPriceResponse, error) {

	if req.Username == "" {
		return &pb.CalculateTotalPriceResponse{
			TotalPrice:   0.0,
			ErrorMessage: "Username must be provided and not empty",
		}, status.Error(codes.InvalidArgument, "Username must be provided and not empty")
	}

	totalPrice, err := s.repo.CalculateTotalPrice(req.Username)
	if err != nil {
		return &pb.CalculateTotalPriceResponse{TotalPrice: 0.0, ErrorMessage: err.Error()}, err
	}
	return &pb.CalculateTotalPriceResponse{TotalPrice: totalPrice}, nil
}
