package handler

import (
	"github.com/gorilla/mux"
)

type Handler struct {
}

func (h *Handler) InitRoute() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/users/", GetUsers)
	router.HandleFunc("/user/{id:[0-9]+}", GetUser)
	router.HandleFunc("/createUser/", SetUser)
	router.HandleFunc("/SignUser/", SignUser)
	router.HandleFunc("/editUsers/", PatchUser)
	router.HandleFunc("/deleteUsers/", DeleteUser)

	return router
}
