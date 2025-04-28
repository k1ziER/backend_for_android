package service

import (
	"android/internal/domain"
	"android/internal/repository"
)

type User interface {
	CreateUser(user domain.User) (int, error)
}

type Service struct {
	User
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		User: NewAuthService(repo.User),
	}
}
