package service

import (
	"android/internal/repository"
	"android/pkg/ports"
)

type Service struct {
	User ports.UserService
}

func NewService(repo *repository.Repository, blackList *UserBlackList) *Service {
	return &Service{
		User: NewAuthService(repo.User, blackList),
	}
}
