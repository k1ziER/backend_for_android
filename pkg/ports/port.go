package ports

import (
	"android/pkg/domain"
)

type UserRepo interface {
	CreateUser(user domain.User) (domain.User, error)
	SignIn(login, password string) (domain.User, error)
	GetUser(id int) (domain.User, error)
	UpdateUser(user domain.User) error
	DeleteUser(id int) error
}

type UserService interface {
	SignIn(login, password string) (domain.User, error)
	CreateUser(user domain.User) (domain.User, error)
	GenerateToken(user domain.User) (string, error)
	ParseToken(token string) (int, error)
	GetUser(id int) (domain.User, error)
	UpdateUser(user domain.User) error
	DeleteUser(id int) error
	Logout(token string)
	CreateTicket(id int, input domain.Ticket) (string, error)
	ParseTicketToken(accessToken string) (domain.Ticket, error)
}

type UserBlackList interface {
	AddUserBlackList(userId int)
	IsUserBlackListed(userId int) bool
	AddTokenBlackList(token string)
	IsTokenBlackListed(token string) bool
}
