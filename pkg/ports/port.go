package ports

import (
	"android/pkg/domain"
	"context"
)

type UserRepo interface {
	CreateUser(user domain.User) (domain.User, error)
	SignIn(ctx context.Context, login, password string) (domain.User, error)
	GetUser(ctx context.Context, id int) (domain.User, error)
	UpdateUser(user domain.User) error
	DeleteUser(id int) error
}

type UserService interface {
	SignIn(ctx context.Context, login, password string) (domain.User, error)
	CreateUser(user domain.User) (domain.User, error)
	GenerateToken(user domain.User) (string, error)
	ParseToken(token string) (int, error)
	GetUser(ctx context.Context, id int) (domain.User, error)
	UpdateUser(user domain.User) error
	DeleteUser(id int) error
	Logout(token string)
	CreateTicket(id int, input domain.Ticket) (string, error)
	ParseTicketToken(accessToken string) (domain.Ticket, error)
	GenerateUUIDString() string
}

type UserBlackList interface {
	AddUserBlackList(userId int)
	IsUserBlackListed(userId int) bool
	AddTokenBlackList(token string)
	IsTokenBlackListed(token string) bool
}
