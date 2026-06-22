package clients

import (
	pbAuth "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/auth"
	pbCart "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/cart"
	pbCatalog "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/catalog"
	pbOrder "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/order"
	pbPayment "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/payment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ServiceClients groups all gRPC clients
type ServiceClients struct {
	Auth        pbAuth.AuthenticationServiceClient
	Cart        pbCart.CartServiceClient
	Catalog     pbCatalog.CatalogServiceClient
	Order       pbOrder.OrderServiceClient
	Payment     pbPayment.PaymentServiceClient
	connections []*grpc.ClientConn
}

// InitClients initializes all gRPC connections
func InitClients() (*ServiceClients, error) {
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())

	// Auth connection
	authConn, err := grpc.NewClient("localhost:8081", opts)
	if err != nil {
		return nil, err
	}

	// Cart connection
	cartConn, err := grpc.NewClient("localhost:8082", opts)
	if err != nil {
		return nil, err
	}

	// Catalog connection
	catalogConn, err := grpc.NewClient("localhost:8083", opts)
	if err != nil {
		return nil, err
	}

	// Order connection
	orderConn, err := grpc.NewClient("localhost:8084", opts)
	if err != nil {
		return nil, err
	}

	// Payment connection
	paymentConn, err := grpc.NewClient("localhost:8085", opts)
	if err != nil {
		return nil, err
	}

	return &ServiceClients{
		Auth:        pbAuth.NewAuthenticationServiceClient(authConn),
		Cart:        pbCart.NewCartServiceClient(cartConn),
		Catalog:     pbCatalog.NewCatalogServiceClient(catalogConn),
		Order:       pbOrder.NewOrderServiceClient(orderConn),
		Payment:     pbPayment.NewPaymentServiceClient(paymentConn),
		connections: []*grpc.ClientConn{authConn, cartConn, catalogConn, orderConn, paymentConn},
	}, nil
}

// Close closes all connections when the server shuts down
func (s *ServiceClients) Close() {
	for _, conn := range s.connections {
		conn.Close()
	}
}
