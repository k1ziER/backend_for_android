package handler

import (
	"android/internal/service"
	"fmt"

	"github.com/gorilla/mux"
)

type Handler struct {
	services *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{services: service}
}

func (h *Handler) InitRoute() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/users/", GetUsers)
	router.HandleFunc("/user/{id:[0-9]+}", h.GetUser)
	router.HandleFunc("/createUser/", h.SetUser)
	router.HandleFunc("/editUsers/", PatchUser)
	router.HandleFunc("/deleteUsers/", DeleteUser)
	fmt.Println("Start server at :8080")
	return router
}
