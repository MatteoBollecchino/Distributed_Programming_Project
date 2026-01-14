package domain

// pb "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto"

type Role string

const (
	UserRole  Role = "USER"
	AdminRole Role = "ADMIN"
)

// User represents a user account in the system.
type User struct {
	// Username is the unique identifier for the user.
	Username string

	// Password is the hashed password for the user.
	Password string

	// Role is what defines the user's permissions
	Role Role
}

/*
// ModelUserToProtoUser converts a model.User into a pb.User
func ModelUserToProtoUser(user *User) (*pb.User, error) {
	if user == nil {
		return nil, fmt.Errorf("Input argument is nil")
	}

	return &pb.User{
		Username: user.Username,
		Role:     user.Role,
	}, nil
}*/
