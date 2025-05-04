package domain

type User struct {
	Id       int    `json:"id" db:"id"`
	UserName string `json: "userName" binding:"required"`
	Surname  string `json: "surname" binding:"required"`
	Email    string `json: "email" binding:"required"`
	Password string `json: "password" binding:"required"`
	IsAdmin  string `json: "isAdmin" binding:"required"`
	Birthday string `json: birthday binding:"required"`
	Age      int    `json: "age" binding:"required"`
}
