package handler

import (
	"context"
)

// func (h *Handler) singUp(w http.ResponseWriter, r *http.Request) {
// 	input := domain.User{}

// 	decoder := json.NewDecoder(r.Body)
// 	if err := decoder.Decode(&input); err != nil {
// 		newErrorResponse(w, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	go func(input domain.User) {
// 		data, err := json.Marshal(&input)
// 		if err != nil {
// 			logrus.Println(err)
// 			return
// 		}

// 		//на месте этой заглушки должна быть кафка
// 		if data == nil {
// 			logrus.Println(data)
// 		}
// 	}(input)

// 	id, err := h.services.CreateUser(input)
// 	if err != nil {
// 		newErrorResponse(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(map[string]interface{}{
// 		"id": id,
// 	})
// }

func (h *Handler) singIn(con context.Context) {

}
