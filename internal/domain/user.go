package user

type User struct {
	Id       int    `json:"id"`
	UserName string `json: "userName"`
	Surname  string `json: "surname"`
	Email    string `json: "email"`
	Password string `json: "password"`
	Age      int    `json: "age"`
	Birthday string `json: birthday`
}
