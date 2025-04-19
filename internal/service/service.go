package service

import "android/internal/repository"

type User interface {
}

type Service struct {
	User
}

func NewService(repo *repository.Repository) *Service {
	return &Service{}
}
