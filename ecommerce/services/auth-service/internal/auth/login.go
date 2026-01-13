package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/services/auth-service/internal/domain"
)

type LoginUseCase struct {
	repo domain.AuthServiceInterface
}

func NewLoginUseCase(repo domain.AuthServiceInterface) *LoginUseCase {
	return &LoginUseCase{repo: repo}
}

func (uc *LoginUseCase) Execute(username, password string) (*domain.User, error) {
	u, err := uc.repo.GetUser(username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return u, nil
}
