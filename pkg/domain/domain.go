package domain

import "time"

type User struct {
	Id       int     `json:"id" db:"id"`
	UserName string  `json:"userName" db:"username" binding:"required"`
	Login    string  `json:"login" db:"loginn" binding:"required"`
	Surname  *string `json:"surname" db:"surname" binding:"required"`
	Email    string  `json:"email" db:"email" binding:"required"`
	Password string  `json:"password" db:"password_hash" binding:"required"`
}

type Ticket struct {
	TitleAttraction string    `json:"titleAttraction"`
	Date            time.Time `json:"date"`
	Count           int       `json:"count"`
}
