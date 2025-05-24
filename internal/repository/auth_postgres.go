package repository

import (
	"android/pkg/domain"
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
	query := fmt.Sprintf("INSERT INTO %s (userName, loginn, email, password_hash) values ($1, $2, $3, $4) RETURNING id", peoplesTable)
	row := r.db.QueryRow(query, user.UserName, user.Login, user.Email, user.Password)
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

func (r *AuthPostgres) UpdateUser(user domain.User) error {
	query := fmt.Sprintf(`UPDATE %s SET userName = $1,
	 surname = $2,
	 email = $3,
	 loginn = $4
	 WHERE id = $5`, peoplesTable)
	_, err := r.db.Exec(query, user.UserName, user.Surname, user.Email, user.Login, user.Id)
	if err != nil {
		logrus.Println(err)
		return err
	}
	return err
}

func (r *AuthPostgres) DeleteUser(id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1 ", peoplesTable)
	_, err := r.db.Exec(query, id)
	if err != nil {
		logrus.Println(err)
	}

	return err
}
