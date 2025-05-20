package domain

type User struct {
	Id       int    `json:"id" db:"id"`
	UserName string `json:"userName" db:"username" binding:"required"`
	Surname  string `json:"surname" db:"surname" binding:"required"`
	Email    string `json:"email" db:"email" binding:"required"`
	Password string `json:"password" db:"password_hash" binding:"required"`
	IsAdmin  bool   `json:"isAdmin" db:"is_admin" binding:"required"`
	Birthday string `json:"birthday" db:"birthday" binding:"required"`
	Age      int    `json:"age" db:"age" binding:"required"`
}
