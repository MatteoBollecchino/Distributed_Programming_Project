package internal

import (
	"context"

	pb "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/catalog"
	"github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/services/catalog-service/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CatalogServer implements the catalog service gRPC server.
type CatalogServer struct {
	pb.CatalogServiceServer
	repo domain.CatalogServiceInterface
}

func NewCatalogServer(repo domain.CatalogServiceInterface) *CatalogServer {
	return &CatalogServer{repo: repo}
}

// AddItemToCart adds an item to the cart of a specific user.
func (s *CatalogServer) AddItemToCart(ctx context.Context, req *pb.AddCatalogItemRequest) (*pb.AddCatalogItemResponse, error) {

	if req.Item.ItemId == "" || req.Item.Description == "" || req.Item == nil {
		return &pb.AddCatalogItemResponse{
			ErrorMessage: "ItemId and Description must be provided and not empty or nil",
		}, status.Error(codes.InvalidArgument, "ItemId and Description must be provided and not empty or nil")
	}

	if req.Item.QuantityAvailable == 0 {
		return &pb.AddCatalogItemResponse{
			ErrorMessage: "Quantity available must be greater than zero",
		}, status.Error(codes.InvalidArgument, "Quantity available must be greater than zero")
	}

	if req.Item.Price < 0 {
		return &pb.AddCatalogItemResponse{
			ErrorMessage: "Price must be non-negative",
		}, status.Error(codes.InvalidArgument, "Price must be non-negative")
	}

	err := s.repo.AddCatalogItem(req.Item)
	if err != nil {
		return &pb.AddCatalogItemResponse{ErrorMessage: err.Error()}, err
	}
	return &pb.AddCatalogItemResponse{}, nil
}

// RemoveCatalogItem removes a catalog item from the catalog by its unique identifier.
func (s *CatalogServer) RemoveCatalogItem(ctx context.Context, req *pb.RemoveCatalogItemRequest) (*pb.RemoveCatalogItemResponse, error) {

	if req.ItemId == "" {
		return &pb.RemoveCatalogItemResponse{
			ErrorMessage: "ItemId must be provided and not empty",
		}, status.Error(codes.InvalidArgument, "ItemId must be provided and not empty")
	}

	err := s.repo.RemoveCatalogItem(req.ItemId)
	if err != nil {
		return &pb.RemoveCatalogItemResponse{ErrorMessage: err.Error()}, err
	}
	return &pb.RemoveCatalogItemResponse{}, nil
}

// GetCatalogItem retrieves a catalog item by its unique identifier.
func (s *CatalogServer) GetCatalogItem(ctx context.Context, req *pb.GetCatalogItemRequest) (*pb.GetCatalogItemResponse, error) {

	if req.ItemId == "" {
		return &pb.GetCatalogItemResponse{
			Item:         nil,
			ErrorMessage: "ItemId must be provided and not empty",
		}, status.Error(codes.InvalidArgument, "ItemId must be provided and not empty")
	}

	item, err := s.repo.GetCatalogItem(req.ItemId)
	if err != nil {
		return &pb.GetCatalogItemResponse{Item: nil, ErrorMessage: err.Error()}, err
	}
	return &pb.GetCatalogItemResponse{Item: item}, nil
}

// UpdateItemQuantity updates the quantity of an item in the cart of a specific user
func (s *CatalogServer) UpdateQuantityAvailable(ctx context.Context, req *pb.UpdateQuantityAvailableRequest) (*pb.UpdateQuantityAvailableResponse, error) {

	if req.ItemId == "" {
		return &pb.UpdateQuantityAvailableResponse{
			ErrorMessage: "ItemId must be provided and not empty",
		}, status.Error(codes.InvalidArgument, "ItemId must be provided and not empty")
	}

	if req.Quantity == 0 {
		return &pb.UpdateQuantityAvailableResponse{
			ErrorMessage: "Quantity must be greater than zero",
		}, status.Error(codes.InvalidArgument, "Quantity must be greater than zero")
	}

	err := s.repo.UpdateQuantityAvailable(req.ItemId, req.Quantity)
	if err != nil {
		return &pb.UpdateQuantityAvailableResponse{ErrorMessage: err.Error()}, err
	}
	return &pb.UpdateQuantityAvailableResponse{}, nil
}

// UpdatePrice updates the price of a catalog item.
func (s *CatalogServer) UpdatePrice(ctx context.Context, req *pb.UpdatePriceRequest) (*pb.UpdatePriceResponse, error) {

	if req.ItemId == "" {
		return &pb.UpdatePriceResponse{
			ErrorMessage: "ItemId must be provided and not empty",
		}, status.Error(codes.InvalidArgument, "ItemId must be provided and not empty")
	}

	if req.Price < 0 {
		return &pb.UpdatePriceResponse{
			ErrorMessage: "Price must be non-negative",
		}, status.Error(codes.InvalidArgument, "Price must be non-negative")
	}

	err := s.repo.UpdatePrice(req.ItemId, req.Price)
	if err != nil {
		return &pb.UpdatePriceResponse{ErrorMessage: err.Error()}, err
	}

	return &pb.UpdatePriceResponse{}, nil
}

// ListCatalogItems retrieves all catalog items.
func (s *CatalogServer) ListCatalogItems(ctx context.Context, req *pb.ListCatalogItemsRequest) (*pb.ListCatalogItemsResponse, error) {

	catalogItems, err := s.repo.ListCatalogItems()
	if err != nil {
		return &pb.ListCatalogItemsResponse{Items: nil, ErrorMessage: err.Error()}, err
	}
	return &pb.ListCatalogItemsResponse{Items: catalogItems}, nil
}
