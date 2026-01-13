/*DA CONTROLLARE*/

/*package internal

import (
	"context"

	pb "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/services/auth-service/proto"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	repo AuthRepositoryInterface
}

func NewAuthServer(repo AuthRepositoryInterface) *AuthServer {
	return &AuthServer{repo: repo}
}

func (s *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := s.repo.Login(req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	return &pb.LoginResponse{Username: user.Username}, nil
}
func (s *AuthServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	err := s.repo.Register(req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	return &pb.RegisterResponse{Message: "User registered successfully"}, nil
}
func (s *AuthServer) GetAllUsers(ctx context.Context, req *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error) {
	users, err := s.repo.GetAllUsers()
	if err != nil {
		return nil, err
	}
	var pbUsers []*pb.User
	for _, user := range users {
		pbUsers = append(pbUsers, &pb.User{Username: user.Username})
	}
	return &pb.GetAllUsersResponse{Users: pbUsers}, nil
}*/
