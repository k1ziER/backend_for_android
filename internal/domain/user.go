package domain

type User struct {
	Id       int    `json:"id"`
	UserName string `json: "userName" binding:"required"`
	Surname  string `json: "surname" binding:"required"`
	Email    string `json: "email" binding:"required"`
	Password string `json: "password" binding:"required"`
	Birthday string `json: birthday binding:"required"`
	Age      int    `json: "age" binding:"required"`
}
