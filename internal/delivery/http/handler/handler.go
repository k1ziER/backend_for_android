package handler

import (
	"android/internal/service"

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
	router.HandleFunc("/signIn/", h.SignIn)
	router.HandleFunc("/createUser/", h.SetUser)

	apiRouter := router.PathPrefix("/api/").Subrouter()
	apiRouter.Use(h.userIdentity)
	apiRouter.HandleFunc("/users/", h.GetUsers)
	apiRouter.HandleFunc("/user/{id:[0-9]+}", h.GetUser)
	apiRouter.HandleFunc("/editUsers/", h.PatchUser)
	apiRouter.HandleFunc("/deleteUsers/", h.DeleteUser)

	return router
}
