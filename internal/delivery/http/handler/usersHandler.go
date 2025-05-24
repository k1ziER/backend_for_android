package handler

import (
	"android/pkg/domain"
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

func (h *Handler) PatchUser(w http.ResponseWriter, r *http.Request) {
	input := domain.User{}
	userId := r.Context().Value(userCtx)
	input.Id, _ = userId.(int)

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err := h.services.User.UpdateUser(input)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	token, err := h.services.User.GenerateToken(input)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": token,
	})
	go func(input domain.User) {
		data, err := json.Marshal(&input)
		if err != nil {
			logrus.Println(err)
			return
		}

		//на месте этой заглушки должна быть кафка
		if data == nil {
			logrus.Println(data)
		}
	}(input)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userCtx)

	id, _ := userId.(int)
	err := h.services.User.DeleteUser(id)
	if err != nil {
		logrus.Println(err)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"DeletedId": id,
	})
}
