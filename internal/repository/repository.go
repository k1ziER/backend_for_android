package repository

import (
	"android/pkg/ports"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	User ports.UserRepo
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User: NewAuthPostgres(db),
	}
}
