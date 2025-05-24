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

func (r *AuthPostgres) CreateUser(user domain.User) (domain.User, error) {
	query := fmt.Sprintf("INSERT INTO %s (userName, loginn, surname, email, password_hash) values ($1, $2, $3, $4, $5) RETURNING id", peoplesTable)
	row := r.db.QueryRow(query, user.UserName, user.Login, user.Surname, user.Email, user.Password)
	err := row.Scan(&user.Id)
	if err != nil {
		logrus.Println(err)
		return user, err
	}
	return user, err
}

func (r *AuthPostgres) SignIn(login, password string) (domain.User, error) {
	var user domain.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE loginn=$1 AND password_hash=$2", peoplesTable)
	err := r.db.Get(&user, query, login, password)
	if err != nil {
		logrus.Println(err)
	}

	return user, err
}

func (r *AuthPostgres) GetUser(id int) (domain.User, error) {
	var user domain.User
	query := fmt.Sprintf("SELECT userName, loginn, surname, email FROM %s WHERE id=$1 ", peoplesTable)
	err := r.db.Get(&user, query, id)
	if err != nil {
		logrus.Println(err)
	}

	return user, err
}
