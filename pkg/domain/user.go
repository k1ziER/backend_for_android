package domain

type User struct {
	Id       int     `json:"id" db:"id"`
	UserName string  `json:"userName" db:"username" binding:"required"`
	Login    string  `json:"login" db:"loginn" binding:"required"`
	Surname  *string `json:"surname" db:"surname" binding:"required"`
	Email    string  `json:"email" db:"email" binding:"required"`
	Password string  `json:"password" db:"password_hash" binding:"required"`
}
