package domain

/*
Repository defines the interface for user data access operations.
*/
type Repository interface {
	FindByUsername(username string) (*User, error)
	Create(user *User) error
}
