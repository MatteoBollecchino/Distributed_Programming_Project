package clients

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"

	pb "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/auth"
)

type AuthGRPCClient struct {
	client pb.AuthenticationServiceClient
}

func NewAuthGRPCClient(addr string) (*AuthGRPCClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
	if err != nil {
		return nil, err
	}

	// Set a timeout for the connection
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return &AuthGRPCClient{client: pb.NewAuthenticationServiceClient(conn)}, nil
}

func (c *AuthGRPCClient) Login(ctx context.Context, username, password string) (*pb.LoginResponse, error) {

	resp, err := c.client.Login(ctx, &pb.LoginRequest{
		Username: username,
		Password: password,
	})

	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.Unauthenticated {
			return nil, errors.New("invalid credentials")
		}
		return nil, errors.New("auth service unavailable")
	}

	return resp, nil
}
