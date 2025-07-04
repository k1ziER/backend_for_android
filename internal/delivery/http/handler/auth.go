package handler

import (
	"android/internal/kafka"
	"android/pkg/domain"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

const (
	topic         = "setUser"
	consumerGroup = "my-consumer-group"
)

// @Summary SignIn
// @Tags auth
// @Description This API sign in account
// @ID sign in
// @Accept json
// @Produce json
// @Param input body domain.User true "account info"
// @Success 200 {integer} integer
// @Failure 400,404 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Failure default {object} handler.ErrorResponse
// @Router /signIn/ [post]
func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	input := domain.User{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	signIn, err := h.services.User.SignIn(ctx, input.Login, input.Password)
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

}

// @Summary SetUser
// @Tags auth
// @Description This API create account
// @ID create account
// @Accept json
// @Produce json
// @Param input body domain.User true "account info"
// @Success 200 {integer} integer
// @Failure 400,404 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Failure default {object} handler.ErrorResponse
// @Router /createUser/ [post]
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
	err = godotenv.Load()
	var adress []string
	adress = append(adress, os.Getenv("kafka1"), os.Getenv("kafka2"), os.Getenv("kafka3"))

	p, err := kafka.NewProducer(adress)
	if err != nil {
		logrus.Fatal(err)
	}
	message := strings.Join([]string{us.Login, us.Password, us.UserName, us.Email}, " ")
	key := h.services.User.GenerateUUIDString()
	err = p.Produce(message, topic, key)
	if err != nil {
		logrus.Println("Error in produce kafka: %w", err.Error())
	}

}
func (h *Handler) SendToKafka(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load()
	var adress []string
	adress = append(adress, os.Getenv("kafka1"), os.Getenv("kafka2"), os.Getenv("kafka3"))
	c, err := kafka.NewConsumer(adress, topic, consumerGroup)
	if err != nil {
		logrus.Fatal(err)
	}
	go func() {
		c.Start()
	}()

	signChan := make(chan os.Signal, 1)
	signal.Notify(signChan, syscall.SIGINT, syscall.SIGTERM)

	<-signChan
}
