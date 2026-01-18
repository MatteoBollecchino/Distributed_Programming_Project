package internal

import (
	"context"
	"strings"

	pb "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/auth"
	"github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/services/auth-service/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	if err := checkCredetials(req.Username, req.Password); err != nil {
		return &pb.LoginResponse{User: nil, ErrorMessage: err.Error()}, err
	}

	user, err := s.repo.Login(req.Username, req.Password)
	if err != nil {
		return &pb.LoginResponse{User: nil, ErrorMessage: err.Error()}, err
	}
	return &pb.LoginResponse{User: &pb.User{Username: user.Username}}, nil
}

// Register creates a new user account with the provided info.
func (s *AuthServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {

	if err := checkCredetials(req.Username, req.Password); err != nil {
		return &pb.RegisterResponse{ErrorMessage: err.Error()}, err
	}

	err := s.repo.Register(req.Username, req.Password)
	if err != nil {
		return &pb.RegisterResponse{ErrorMessage: err.Error()}, err
	}
	return &pb.RegisterResponse{}, nil
}

// ChangePassword updates the password for the specified user.
func (s *AuthServer) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {

	if err := checkCredetials(req.Username, req.NewPassword); err != nil {
		return &pb.ChangePasswordResponse{ErrorMessage: err.Error()}, err
	}

	err := s.repo.ChangePassword(req.Username, req.OldPassword, req.NewPassword)
	if err != nil {
		return &pb.ChangePasswordResponse{ErrorMessage: err.Error()}, err
	}
	return &pb.ChangePasswordResponse{}, nil
}

// GetUser retrieves the user information for the specified username.
func (s *AuthServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {

	if err := validUsername(req.Username); err != nil {
		return &pb.GetUserResponse{User: nil, ErrorMessage: err.Error()}, err
	}

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

// PRIVATE FUNCTIONS TO VALIDATE INPUTS

// checkCredetials validates the provided username and password.
func checkCredetials(username, password string) error {
	if err := validUsername(username); err != nil {
		return err
	}
	if err := validPassword(password); err != nil {
		return err
	}
	return nil
}

// validUsername checks if the username is not an empty line.
func validUsername(username string) error {

	if username == "" {
		return status.Error(codes.InvalidArgument, "Username must be provided and not empty")
	}

	return nil
}

// validPassword checks if the password has at least 8 characters, contains a number, and a special character.
func validPassword(password string) error {

	if len(password) < 8 {
		return status.Error(codes.InvalidArgument, "Password must be at least 8 characters long")
	}

	hasNumber := false
	hasSpecial := false
	specialCharacters := "!@#$%^&*()-+"

	for _, char := range password {
		if char >= '0' && char <= '9' {
			hasNumber = true
		}
		if strings.Contains(specialCharacters, string(char)) {
			hasSpecial = true
		}
	}

	if !hasNumber {
		return status.Error(codes.InvalidArgument, "Password must contain at least one number")
	}

	if !hasSpecial {
		return status.Error(codes.InvalidArgument, "Password must contain at least one special character")
	}

	return nil
}
