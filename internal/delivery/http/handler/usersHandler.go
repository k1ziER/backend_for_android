package handler

import (
	"android/internal/kafka"
	"android/pkg/domain"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var countTicket int

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// @Summary GetUser
// @Security ApiKeyAuth
// @Tags api
// @Description This API get account
// @ID get account
// @Accept json
// @Produce json
// @Success 200 {integer} integer
// @Failure 400,404 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Failure default {object} handler.ErrorResponse
// @Router /api/getUser/ [get]
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userCtx)
	ctx := r.Context()

	id, _ := userId.(int)
	user, err := h.services.User.GetUser(ctx, id)
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

// @Summary PatchUser
// @Security ApiKeyAuth
// @Tags api
// @Description This API edit account
// @ID edit account
// @Accept json
// @Produce json
// @Param input body domain.User true "account info"
// @Success 200 {integer} integer
// @Failure 400,404 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Failure default {object} handler.ErrorResponse
// @Router /api/editUser/ [patch]
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
	err = godotenv.Load()
	var adress []string
	adress = append(adress, os.Getenv("kafka1"), os.Getenv("kafka2"), os.Getenv("kafka3"))

	p, err := kafka.NewProducer(adress)
	if err != nil {
		logrus.Fatal(err)
	}
	message := strings.Join([]string{input.Login, input.Password, input.UserName, *input.Surname, input.Email}, " ")
	key := h.services.User.GenerateUUIDString()
	err = p.Produce(message, topic, key)
	if err != nil {
		logrus.Println("Error in produce kafka: %w", err.Error())
	}
}

// @Summary DeleteUser
// @Security ApiKeyAuth
// @Tags api
// @Description This API delete account
// @ID delete account
// @Accept json
// @Produce json
// @Success 200 {integer} integer
// @Failure 400,404 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Failure default {object} handler.ErrorResponse
// @Router /api/deleteUser/ [delete]
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userCtx)

	id, _ := userId.(int)
	err := h.services.User.DeleteUser(id)
	if err != nil {
		logrus.Println(err)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"deletedId": id,
	})
	err = godotenv.Load()
	var adress []string
	adress = append(adress, os.Getenv("kafka1"), os.Getenv("kafka2"), os.Getenv("kafka3"))

	p, err := kafka.NewProducer(adress)
	if err != nil {
		logrus.Fatal(err)
	}
	message := fmt.Sprintf("deleted user: %v", id)
	key := h.services.User.GenerateUUIDString()
	err = p.Produce(message, topic, key)
	if err != nil {
		logrus.Println("Error in produce kafka: %w", err.Error())
	}
}

// @Summary Logout
// @Security ApiKeyAuth
// @Tags api
// @Description This API logout account
// @ID logout account
// @Accept json
// @Produce json
// @Success 200 {integer} integer
// @Failure 400,404 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Failure default {object} handler.ErrorResponse
// @Router /api/logout/ [post]
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userCtx)
	id, _ := userId.(int)

	header := r.Header.Get(authorizationHeader)
	if header == "" {
		newErrorResponse(w, http.StatusUnauthorized, "empty auth header")
		return
	}
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(w, http.StatusUnauthorized, "invalid auth header")
		return
	}

	token := headerParts[1]
	h.services.User.Logout(token)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"logout": id,
	})
}

// @Summary CreateTicket
// @Security ApiKeyAuth
// @Tags api
// @Description This API create ticket
// @ID create ticket
// @Accept json
// @Produce json
// @Param input body domain.User true "account info"
// @Success 200 {integer} integer
// @Failure 400,404 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Failure default {object} handler.ErrorResponse
// @Router /api/createTicket/ [post]
func (h *Handler) CreateTicket(w http.ResponseWriter, r *http.Request) {

	countTicket++
	userId := r.Context().Value(userCtx)

	id, _ := userId.(int)

	input := domain.Ticket{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.services.User.CreateTicket(id, input)
	if err != nil {
		logrus.Println(err)
	}

	ticket, err := h.services.User.ParseTicketToken(token)
	if err != nil {
		logrus.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"titleAttraction": ticket.TitleAttraction,
		"countTicket":     countTicket,
	})
}

// @Summary WebSocket
// @Tags api
// @Description This API connect to chat
// @ID connect chat
// @Accept json
// @Produce json
// @Success 200 {integer} integer
// @Failure 400,404 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Failure default {object} handler.ErrorResponse
// @Router /api/websocket/ [post]
func (h *Handler) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Println(err)
		return
	}
	defer func() error {
		err = conn.Close()
		if err != nil {
			logrus.Println(err)
			return err
		}
		return err
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			logrus.Println(err)
			return
		}

		logrus.Printf("Received message: %s", message)

		err = conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			logrus.Println(err)
			return
		}
	}
}
