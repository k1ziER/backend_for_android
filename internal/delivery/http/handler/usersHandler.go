package handler

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userCtx)

	id, _ := userId.(int)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id": id,
	})
}
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) PatchUser(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {

}
