package handler

import (
	"android/internal/domain"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

type signInInput struct {
	Email    string `json: "email" binding:"required"`
	Password string `json: "password" binding:"required"`
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	input := signInInput{}
	//logrus.Println("123")
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.services.User.GenerateToken(input.Email, input.Password)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": token,
	})

	// go func(input signInInput) {
	// 	data, err := json.Marshal(&input)
	// 	if err != nil {
	// 		logrus.Println(err)
	// 		return
	// 	}

	// 	//на месте этой заглушки должна быть кафка
	// 	if data == nil {
	// 		logrus.Println(data)
	// 	}
	// }(input)

}

func (h *Handler) SetUser(w http.ResponseWriter, r *http.Request) {
	input := domain.User{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.User.CreateUser(input)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id": id,
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
