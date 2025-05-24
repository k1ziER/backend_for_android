package service

import (
	"android/internal/domain"
	"android/internal/repository"
)

type User interface {
	SignIn(login, password string) (domain.User, error)
	CreateUser(user domain.User) (domain.User, error)
	GenerateToken(user domain.User) (string, error)
	ParseToken(token string) (int, error)
	GetUser(id int) (domain.User, error)
}

type Service struct {
	User
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		User: NewAuthService(repo.User),
	}
}
