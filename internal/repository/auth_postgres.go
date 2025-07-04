package repository

import (
	"android/internal/redis"
	"android/pkg/domain"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type AuthPostgres struct {
	redis *redis.Client
	db    *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB, redis *redis.Client) *AuthPostgres {
	return &AuthPostgres{
		db:    db,
		redis: redis,
	}
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

func (r *AuthPostgres) SignIn(ctx context.Context, login, password string) (domain.User, error) {
	var user domain.User

	keyCache := fmt.Sprintf("login: %s, password: %s", login, password)

	cachedData, err := r.redis.Get(ctx, keyCache)
	if err == nil {
		err = json.Unmarshal(cachedData, &user)
		if err == nil {
			return user, nil
		}
	}

	query := fmt.Sprintf("SELECT id FROM %s WHERE loginn=$1 AND password_hash=$2", peoplesTable)
	err = r.db.Get(&user, query, login, password)
	if err != nil {
		logrus.Println(err)
	}

	if data, err := json.Marshal(user); err == nil {
		_ = r.redis.Set(ctx, keyCache, data, time.Hour)
	}

	return user, err
}

func (r *AuthPostgres) GetUser(ctx context.Context, id int) (domain.User, error) {

	var user domain.User

	keyCache := fmt.Sprintf("userId: %d", id)

	cachedData, err := r.redis.Get(ctx, keyCache)
	if err == nil {
		err = json.Unmarshal(cachedData, &user)
		if err == nil {
			return user, nil
		}
	}

	query := fmt.Sprintf("SELECT userName, loginn, surname, email FROM %s WHERE id=$1 ", peoplesTable)
	err = r.db.Get(&user, query, id)
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
