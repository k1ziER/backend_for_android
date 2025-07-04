package repository

import (
	"android/internal/redis"
	"android/pkg/ports"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	User ports.UserRepo
}

func NewRepository(db *sqlx.DB, redis *redis.Client) *Repository {
	return &Repository{
		User: NewAuthPostgres(db, redis),
	}
}
