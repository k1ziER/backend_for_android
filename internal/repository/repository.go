package repository

import (
	"android/internal/domain"

	"github.com/jmoiron/sqlx"
)

type User interface {
	CreateUser(user domain.User) (domain.User, error)
	SignIn(login, password string) (domain.User, error)
	GetUser(id int) (domain.User, error)
}

type Repository struct {
	User
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User: NewAuthPostgres(db),
	}
}
