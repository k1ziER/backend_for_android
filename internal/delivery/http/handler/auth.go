package handler

import (
	"android/internal/domain"
	"context"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

func (h *Handler) singUp(w http.ResponseWriter, r *http.Request) {
	input := domain.User{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

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

func (h *Handler) singIn(con context.Context) {

}
