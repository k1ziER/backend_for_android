package handler

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userCtx)

	id, _ := userId.(int)
	user, err := h.services.User.GetUser(id)
	if err != nil {
		logrus.Println(err)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":       id,
		"userName": user.UserName,
		"login":    user.Login,
		"surname":  user.Surname,
		"email":    user.Email,
	})
}

func (h *Handler) PutUser(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userCtx)

	id, _ := userId.(int)
	user, err := h.services.User.GetUser(id)
	if err != nil {
		logrus.Println(err)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":       id,
		"userName": user.UserName,
		"surname":  user.Surname,
		"email":    user.Email,
	})
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {

}
