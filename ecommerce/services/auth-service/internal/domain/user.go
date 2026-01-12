package domain

type Role string

const (
	UserRole  Role = "USER"
	AdminRole Role = "ADMIN"
)

/*
User represents a user account in the system.

Every user is represented in the system by a unique username,
a hashed password for authentication,
and a role that defines the user's permissions.
*/
type User struct {
	Username string
	Password string
	Role     Role
}
