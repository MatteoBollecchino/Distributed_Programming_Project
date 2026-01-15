package tests

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	pb "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/auth"
	"github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/services/auth-service/internal/domain"
	"github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/services/auth-service/internal/repository"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect database: %v", err)
	}

	err = db.AutoMigrate(&domain.User{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}
	return db
}

func setupDefaultUser(t *testing.T, db *gorm.DB, repo *repository.AuthRepository) {
	username := "user1"
	password := "Password1+"

	hashedPassword, err := repo.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	user := domain.User{Username: username, Password: hashedPassword, Role: "USER"}
	err = db.Create(&user).Error
	if err != nil {
		t.Fatalf("Failed to create default user: %v", err)
	}

	if err := db.Where("username = ? AND password = ?", username, hashedPassword).First(&user).Error; err != nil {
		t.Fatalf("Failed to retrieve default user: %v", err)
	}
}

func setupTest(t *testing.T) (*gorm.DB, *repository.AuthRepository) {
	db := setupTestDB(t)
	repo := repository.NewAuthRepository(db)

	setupDefaultUser(t, db, repo)

	return db, repo
}

func TestLoginWithCorrectCredentials(t *testing.T) {
	_, repo := setupTest(t)

	username := "user1"
	password := "Password1+"

	var user *pb.User

	// Now, attempt to login with the same credentials
	user, err := repo.Login(username, password)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if user.Username != username {
		t.Fatalf("Expected username %s, got %s", username, user.Username)
	}
	if user.Role != pb.Role_USER {
		t.Fatalf("Expected role %v, got %v", pb.Role_USER, user.Role)
	}

}

func TestLoginWithIncorrectPassword(t *testing.T) {
	_, repo := setupTest(t)

	username := "user1"
	password := "Prongpassword7+"

	// Now, attempt to login with incorrect credentials
	_, err := repo.Login(username, password)
	if err == nil {
		t.Fatalf("Expected login to fail with incorrect credentials")
	}
}

func TestLoginWithIncorrectUsername(t *testing.T) {
	_, repo := setupTest(t)

	username := "user" // Incorrect username
	password := "Password1+"

	// Attempt to login with invalid username
	_, err := repo.Login(username, password)
	if err == nil {
		t.Fatalf("Expected login to fail with invalid username")
	}
}

func TestLoginWithNonExistentUser(t *testing.T) {
	_, repo := setupTest(t)

	username := "nonexistent"
	password := "SomePassword1+"

	// Attempt to login with non-existent user
	_, err := repo.Login(username, password)
	if err == nil {
		t.Fatalf("Expected login to fail with non-existent user")
	}
}

func TestRegisterNewUser(t *testing.T) {
	db, repo := setupTest(t)

	username := "newuser"
	password := "NewPassword1+"

	// Attempt to register a new user
	err := repo.Register(username, password)
	if err != nil {
		t.Fatalf("Registration failed: %v", err)
	}

	// Verify the user was created
	var user domain.User
	err = db.Where("username = ?", username).First(&user).Error
	if err != nil {
		t.Fatalf("Failed to retrieve registered user: %v", err)
	}
	if user.Username != username {
		t.Fatalf("Expected username %s, got %s", username, user.Username)
	}
	if user.Role != domain.UserRole {
		t.Fatalf("Expected role %s, got %s", domain.UserRole, user.Role)
	}
}

func TestRegisterDuplicateUsername(t *testing.T) {
	_, repo := setupTest(t)

	username := "user1" // Duplicate username
	password := "AnotherPassword1+"

	// Attempt to register with a duplicate username
	err := repo.Register(username, password)
	if err == nil {
		t.Fatalf("Expected registration to fail with duplicate username")
	}
}

func TestRegisterInvalidPassword(t *testing.T) {
	_, repo := setupTest(t)

	username := "validuser"
	password := "short" // Invalid password

	// Attempt to register with an invalid password
	err := repo.Register(username, password)
	if err == nil {
		t.Fatalf("Expected registration to fail with invalid password")
	}
}

func TestRegisterInvalidUsername(t *testing.T) {
	_, repo := setupTest(t)

	username := "" // Invalid username
	password := "ValidPassword1+"

	// Attempt to register with an invalid username
	err := repo.Register(username, password)
	if err == nil {
		t.Fatalf("Expected registration to fail with invalid username")
	}
}

func TestRegisterAndLoginFlow(t *testing.T) {
	_, repo := setupTest(t)

	username := "flowuser"
	password := "FlowPassword1+"

	// Register the user
	err := repo.Register(username, password)
	if err != nil {
		t.Fatalf("Registration failed: %v", err)
	}

	// Now, attempt to login with the same credentials
	user, err := repo.Login(username, password)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if user.Username != username {
		t.Fatalf("Expected username %s, got %s", username, user.Username)
	}
}

func TestChangePasswordSuccess(t *testing.T) {
	_, repo := setupTest(t)

	username := "user1"
	oldPassword := "Password1+"
	newPassword := "NewPassword1+"

	// Attempt to change password
	err := repo.ChangePassword(username, oldPassword, newPassword)
	if err != nil {
		t.Fatalf("ChangePassword failed: %v", err)
	}

	// Verify that login with new password works
	_, err = repo.Login(username, newPassword)
	if err != nil {
		t.Fatalf("Login with new password failed: %v", err)
	}

	// Verify that login with old password fails
	_, err = repo.Login(username, oldPassword)
	if err == nil {
		t.Fatalf("Login with old password should have failed")
	}
}

func TestGetCorrectUser(t *testing.T) {
	_, repo := setupTest(t)

	username := "user1"
	user, err := repo.GetUser(username)
	if err != nil {
		t.Fatalf("GetUser failed: %v", err)
	}
	if user.Username != username {
		t.Fatalf("Expected username %s, got %s", username, user.Username)
	}
	if user.Role != pb.Role_USER {
		t.Fatalf("Expected role %s, got %s", pb.Role_USER, user.Role)
	}

}

func TestGetNonExistentUser(t *testing.T) {
	_, repo := setupTest(t)

	username := "nonexistent"
	_, err := repo.GetUser(username)
	if err == nil {
		t.Fatalf("Expected GetUser to fail for non-existent user")
	}
}

func TestGetAllUsers(t *testing.T) {
	db, repo := setupTest(t)

	// list of domain.User to add
	users := []domain.User{}
	users = append(users, domain.User{Username: "mario", Password: "MarioPassword1+"})
	users = append(users, domain.User{Username: "luigi", Password: "LuigiPassword1+"})
	users = append(users, domain.User{Username: "matteo", Password: "MatteoPassword1+"})

	for _, u := range users {
		password, err := repo.HashPassword(u.Password)
		if err != nil {
			t.Fatalf("Failed to hash password: %v", err)
		}

		user := domain.User{Username: u.Username, Password: password, Role: domain.UserRole}
		err = db.Create(&user).Error
		if err != nil {
			t.Fatalf("Failed to create user %s: %v", u.Username, err)
		}
	}

	list, err := repo.GetAllUsers()
	if err != nil {
		t.Fatalf("GetAllUsers failed: %v", err)
	}
	if len(list) != 4 { // 3 added + 1 default
		t.Fatalf("Expected 4 users, got %d", len(list))
	}

	// Verify that all users are returned
	for _, user := range list {
		if user.Username == "mario" || user.Username == "luigi" || user.Username == "matteo" {
			continue
		}
		if user.Username != "user1" {
			t.Fatalf("Unexpected username in list: %s", user.Username)
		}
	}
}

func TestHashPassword(t *testing.T) {
	_, repo := setupTest(t)

	password := "SomePassword1+"
	hashedPassword, err := repo.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	if repo.CheckPassword(hashedPassword, password) != nil {
		t.Fatalf("Hashed password does not match original password: %v", err)
	}
}

func TestCreateAdmin(t *testing.T) {
	db, repo := setupTest(t)
	username := "adminUser"
	password := "AdminPassword1+"

	// Attempt to create an admin user
	err := repo.CreateAdmin(username, password)
	if err != nil {
		t.Fatalf("CreateAdmin failed: %v", err)
	}

	// Verify the admin user was created
	var user domain.User
	err = db.Where("username = ?", username).First(&user).Error
	if err != nil {
		t.Fatalf("Failed to retrieve admin user: %v", err)
	}
	if user.Username != username {
		t.Fatalf("Expected username %s, got %s", username, user.Username)
	}
	if user.Role != "ADMIN" {
		t.Fatalf("Expected role ADMIN, got %s", user.Role)
	}
}

func TestCreateDefaultUsersAdmins(t *testing.T) {
	db, repo := setupTest(t)

	err := repo.CreateDefaultUsersAdmins()
	if err != nil {
		t.Fatalf("CreateDefaultUsersAdmins failed: %v", err)
	}

	users, err := repo.GetAllUsers()
	if err != nil {
		t.Fatalf("GetAllUsers failed: %v", err)
	}

	if len(users) != 5 { // 3 default users + 2 default admins
		t.Fatalf("Expected 5 users, got %d", len(users))
	}

	for _, dUser := range users {
		// Verify the user was created
		var user domain.User
		err = db.Where("username = ?", dUser.Username).First(&user).Error
		if err != nil {
			t.Fatalf("Failed to retrieve registered user: %v", err)
		}
		if user.Username != dUser.Username {
			t.Fatalf("Expected username %s, got %s", dUser.Username, user.Username)
		}
		if domain.Role(dUser.Role.String()) != user.Role {
			t.Fatalf("Expected role %s, got %s", dUser.Role.String(), user.Role)
		}
	}
}
