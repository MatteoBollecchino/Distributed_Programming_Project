package repository

import (
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	pb "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/auth"
	"github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/services/auth-service/internal/domain"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// Login authenticates a user with the given username and password, after validating the credentials.
func (r *AuthRepository) Login(username, password string) (*pb.User, error) {

	if err := checkCredetials(username, password); err != nil {
		return nil, err
	}

	user, err := r.getUserByUserame(username)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("Invalid credentials")
	}

	pbUser, err := domain.DomainUserToProtoUser(user)
	if err != nil {
		return nil, err
	}

	return pbUser, nil
}

// Register creates a new USER account with the provided info, after validating the credentials.
func (r *AuthRepository) Register(username, password string) error {

	if err := checkCredetials(username, password); err != nil {
		return err
	}

	err := uniqueUsername(r.db, username)
	if err != nil {
		return err
	}

	// Hash the password before storing it in the database.
	hashedPassword, err := r.HashPassword(password)
	if err != nil {
		return err
	}

	// Register creates only users with USER role
	user := domain.User{Username: username, Password: hashedPassword, Role: domain.UserRole}
	return r.db.Create(&user).Error
}

// ChangePassword updates the password for a given user, after validating the old credentials and the new password.
func (r *AuthRepository) ChangePassword(username, oldPassword, newPassword string) error {

	checkCredetials(username, oldPassword)

	if err := validPassword(newPassword); err != nil {
		return err
	}

	user, err := r.getUserByUserame(username)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("Invalid password")
	}

	user.Password, err = r.HashPassword(newPassword)
	if err != nil {
		return err
	}

	return r.db.Model(&domain.User{}).Where("username = ?", username).Updates(user).Error
}

// GetUser retrieves the user information for the specified username, after validating the username.
func (r *AuthRepository) GetUser(username string) (*pb.User, error) {

	if err := validUsername(username); err != nil {
		return nil, err
	}

	var user domain.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	pbUser, err := domain.DomainUserToProtoUser(&user)
	if err != nil {
		return nil, err
	}

	return pbUser, nil
}

// GetAllUsers retrieves all users registered in the system.
func (r *AuthRepository) GetAllUsers() ([]*pb.User, error) {
	var users []*domain.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}

	var pbUsers []*pb.User
	for _, user := range users {
		pbUser, err := domain.DomainUserToProtoUser(user)
		if err != nil {
			return nil, err
		}
		pbUsers = append(pbUsers, pbUser)
	}

	return pbUsers, nil
}

// CreateAdmin creates a new ADMIN account with the provided info, after validating the credentials.
// Admins are created separately from regular users.
// (The username must be unique as well)
func (r *AuthRepository) CreateAdmin(username, password string) error {

	if err := checkCredetials(username, password); err != nil {
		return err
	}
	err := uniqueUsername(r.db, username)
	if err != nil {
		return err
	}

	// Hash the password before storing it in the database.
	hashedPassword, err := r.HashPassword(password)
	if err != nil {
		return err
	}
	user := domain.User{Username: username, Password: hashedPassword, Role: domain.AdminRole}
	return r.db.Create(&user).Error
}

// HashPassword hashes the password using bcrypt.
func (r *AuthRepository) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

// CreateDefaultUsers creates default admins and users.
func (r *AuthRepository) CreateDefaultUsersAdmins() error {

	defaultUsers := []domain.User{
		{Username: "Marco", Password: "DefaultPassword1+"},
		{Username: "Noemi", Password: "DefaultPassword2+"},
	}

	defaultAdmins := []domain.User{
		{Username: "adminBolle", Password: "DefaultPassword1+"},
		{Username: "adminDani", Password: "DefaultPassword2+"},
	}

	for _, du := range defaultUsers {
		if err := r.Register(du.Username, du.Password); err != nil {
			return err
		}
	}

	for _, da := range defaultAdmins {
		if err := r.CreateAdmin(da.Username, da.Password); err != nil {
			return err
		}
	}

	return nil
}

// PRIVATE FUNCTIONS AND METHODS

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
		return errors.New("Invalid Username")
	}

	return nil
}

// validPassword checks if the password has at least 8 characters, contains a number, and a special character.
func validPassword(password string) error {

	if len(password) < 8 {
		return errors.New("Password must be at least 8 characters long")
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
		return errors.New("Password must contain at least one number")
	}

	if !hasSpecial {
		return errors.New("Password must contain at least one special character")
	}

	return nil
}

// uniqueUsername checks if the username is unique in the database.
func uniqueUsername(db *gorm.DB, username string) error {
	var user domain.User
	if err := db.Where("username = ?", username).First(&user).Error; err == nil {
		return errors.New("Username already exists")
	}
	return nil
}

// getUserByUserame retrieves a user by username from the database.
// (will be differente from GetUser as it has *pb.User return type)
func (r *AuthRepository) getUserByUserame(username string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
