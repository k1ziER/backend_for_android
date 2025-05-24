package handler

import (
	"android/pkg/domain"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	input := domain.User{}
	//logrus.Println("123")
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	signIn, err := h.services.User.SignIn(input.Login, input.Password)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	token, err := h.services.User.GenerateToken(signIn)
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

	us, err := h.services.User.CreateUser(input)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	token, err := h.services.User.GenerateToken(us)
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
