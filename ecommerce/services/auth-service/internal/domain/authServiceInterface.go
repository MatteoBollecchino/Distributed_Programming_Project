package domain

/*
AuthServiceInterface defines the interface for user data access operations.
*/
type AuthServiceInterface interface {
	// Login authenticates a user with the given username and password.
	Login(username, password string) (*User, error)

	// Register creates a new user account with the provided info.
	Register(username, password string) error

	// ChangePassword updates the password for a given user.
	ChangePassword(username, oldPassword, newPassword string) error

	// GetUser retrieves the user information for the specified username.
	GetUser(username string) (*User, error)

	// GetUsers retrieves all users registered in the system.
	GetUsers() ([]*User, error)
}
