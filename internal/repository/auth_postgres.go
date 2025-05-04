package repository

import (
	"android/internal/domain"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user domain.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (userName, surname, email, password_hash, is_admin, birthday, age) values ($1, $2, $3, $4, $5, $6, $7) RETURNING id", peoplesTable)
	row := r.db.QueryRow(query, user.UserName, user.Surname, user.Email, user.Password, user.IsAdmin, user.Birthday, user.Age)
	err := row.Scan(&id)
	if err != nil {
		logrus.Println(err)
		return 0, err
	}
	return id, err
}

func (r *AuthPostgres) GetUser(login, password string) (domain.User, error) {
	var user domain.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE email=$1 AND password_hash=$2", peoplesTable)
	err := r.db.Get(&user, query, login, password)
	if err != nil {
		logrus.Println(err)
	}

	return user, err
}
