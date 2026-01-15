package domain

import pb "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/auth"

// AuthServiceInterface defines the interface for user data access operations.
type AuthServiceInterface interface {
	// Login authenticates a user with the given username and password.
	Login(username, password string) (*pb.User, error)

	// Register creates a new user account with the provided info.
	Register(username, password string) error

	// ChangePassword updates the password for a given user.
	ChangePassword(username, oldPassword, newPassword string) error

	// GetUser retrieves the user information for the specified username.
	GetUser(username string) (*pb.User, error)

	// GetAllUsers retrieves all users registered in the system.
	GetAllUsers() ([]*pb.User, error)
}
