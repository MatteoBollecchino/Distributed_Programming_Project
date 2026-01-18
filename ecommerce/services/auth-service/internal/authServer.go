package internal

import (
	"context"

	pb "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/auth"
	"github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/services/auth-service/internal/domain"
)

// AuthServer implements the authentication service gRPC server.
type AuthServer struct {
	pb.AuthenticationServiceServer
	repo domain.AuthServiceInterface
}

func NewAuthServer(repo domain.AuthServiceInterface) *AuthServer {
	return &AuthServer{repo: repo}
}

// Login authenticates a user with the given username and password.
func (s *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := s.repo.Login(req.Username, req.Password)
	if err != nil {
		return &pb.LoginResponse{User: nil, ErrorMessage: err.Error()}, err
	}
	return &pb.LoginResponse{User: &pb.User{Username: user.Username}}, nil
}

// Register creates a new user account with the provided info.
func (s *AuthServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	err := s.repo.Register(req.Username, req.Password)
	if err != nil {
		return &pb.RegisterResponse{ErrorMessage: err.Error()}, err
	}
	return &pb.RegisterResponse{}, nil
}

// ChangePassword updates the password for the specified user.
func (s *AuthServer) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	err := s.repo.ChangePassword(req.Username, req.OldPassword, req.NewPassword)
	if err != nil {
		return &pb.ChangePasswordResponse{ErrorMessage: err.Error()}, err
	}
	return &pb.ChangePasswordResponse{}, nil
}

// GetUser retrieves the user information for the specified username.
func (s *AuthServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := s.repo.GetUser(req.Username)
	if err != nil {
		return &pb.GetUserResponse{User: nil, ErrorMessage: err.Error()}, err
	}
	return &pb.GetUserResponse{User: &pb.User{Username: user.Username}}, nil
}

// GetAllUsers retrieves all users registered in the system.
func (s *AuthServer) GetAllUsers(ctx context.Context, req *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error) {
	users, err := s.repo.GetAllUsers()
	if err != nil {
		return &pb.GetAllUsersResponse{Users: nil, ErrorMessage: err.Error()}, err
	}
	return &pb.GetAllUsersResponse{Users: users}, nil
}
