package domain

type Role string

const (
	UserRole  Role = "USER"
	AdminRole Role = "ADMIN"
)

/*
User represents a user account in the system.
*/
type User struct {
	// Username is the unique identifier for the user.
	Username string

	// Password is the hashed password for the user.
	Password string

	// Role is what defines the user's permissions
	Role Role
}
