package handler

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

type error struct {
	Message string `json:"message"`
}

func newErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	logrus.Error(message)

	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(error{Message: message})
	if err != nil {
		logrus.Println(err)
		return
	}
}
