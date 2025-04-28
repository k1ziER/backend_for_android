package service

import (
	"android/internal/domain"
	"android/internal/repository"
	"crypto/sha1"
	"fmt"
)

const solt = "earagdrd4153426dsfdg"

type AuthService struct {
	repo repository.User
}

func NewAuthService(repo repository.User) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user domain.User) (int, error) {
	user.Password = s.generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(solt)))
}
