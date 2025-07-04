package handler

import (
	"android/internal/service"

	_ "android/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handler struct {
	services  *service.Service
	blackList *service.UserBlackList
}

func NewHandler(service *service.Service, blackList *service.UserBlackList) *Handler {
	return &Handler{
		services:  service,
		blackList: blackList,
	}
}

func (h *Handler) InitRoute() *mux.Router {
	router := mux.NewRouter()
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8090/swagger/doc.json"), // URL до вашего swagger.json
	))
	router.HandleFunc("/signIn/", h.SignIn)
	router.HandleFunc("/createUser/", h.SetUser)
	router.HandleFunc("/sendToKafka", h.SendToKafka)

	apiRouter := router.PathPrefix("/api/").Subrouter()
	apiRouter.Use(h.userIdentity)
	apiRouter.HandleFunc("/getUser/", h.GetUser)
	apiRouter.HandleFunc("/editUser/", h.PatchUser)
	apiRouter.HandleFunc("/deleteUser/", h.DeleteUser)
	apiRouter.HandleFunc("/logout/", h.Logout)
	apiRouter.HandleFunc("/createTicket/", h.CreateTicket)
	apiRouter.HandleFunc("/websocket/", h.WebSocketHandler)

	return router
}
