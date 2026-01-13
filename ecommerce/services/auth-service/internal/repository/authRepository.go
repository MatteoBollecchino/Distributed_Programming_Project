package repository

import (
	"github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/services/auth-service/internal/domain"
	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *authRepository {
	return &authRepository{db: db}
}

// Login authenticates a user with the given username and password, after validating the credentials.
func (r *authRepository) Login(username, password string) (*domain.User, error) {

	checkCredetials(username, password)

	var user domain.User
	if err := r.db.Where("username = ? AND password = ?", username, password).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Register creates a new user account with the provided info, after validating the credentials.
func (r *authRepository) Register(username, password string) error {

	checkCredetials(username, password)

	user := domain.User{Username: username, Password: password}
	return r.db.Create(&user).Error
}

// ChangePassword updates the password for a given user, after validating the old credentials and the new password.
func (r *authRepository) ChangePassword(username, oldPassword, newPassword string) error {

	checkCredetials(username, oldPassword)

	if err := validPassword(newPassword); err != nil {
		return err
	}

	var user domain.User
	if err := r.db.Where("username = ? AND password = ?", username, oldPassword).First(&user).Error; err != nil {
		return err
	}
	user.Password = newPassword
	return r.db.Save(&user).Error
}

// GetUser retrieves the user information for the specified username, after validating the username.
func (r *authRepository) GetUser(username string) (*domain.User, error) {

	if err := validUsername(username); err != nil {
		return nil, err
	}

	var user domain.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUsers retrieves all users registered in the system.
func (r *authRepository) GetUsers() ([]*domain.User, error) {
	var users []*domain.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

/* Inserire le funzioni per l'hashing delle password degli utenti */

func checkCredetials(username, password string) error {
	if err := validUsername(username); err != nil {
		return err
	}
	if err := validPassword(password); err != nil {
		return err
	}
	return nil
}

func validUsername(username string) error {
	// Implement username validation logic here
	return nil
}

func validPassword(password string) error {
	// Implement password validation logic here
	return nil
}
