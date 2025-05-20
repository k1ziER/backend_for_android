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
		"surname":  user.Surname,
		"email":    user.Email,
		"isAdmin":  user.IsAdmin,
		"birthday": user.Birthday,
		"age":      user.Age,
	})
}
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userCtx)

	id, _ := userId.(int)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) PatchUser(w http.ResponseWriter, r *http.Request) {
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
		"isAdmin":  user.IsAdmin,
		"birthday": user.Birthday,
		"age":      user.Age,
	})
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {

}
