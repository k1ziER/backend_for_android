package handler

import (
	"context"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) userIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//fmt.Println("Все заголовки:", r.Header)
		header := r.Header.Get(authorizationHeader)

		if header == "" {
			newErrorResponse(w, http.StatusUnauthorized, "empty auth header")
			return
		}

		headerParts := strings.Replace(header, "Bearer ", "", 1)
		blackToken := h.blackList.IsTokenBlackListed(headerParts)
		if blackToken == true {
			newErrorResponse(w, http.StatusUnauthorized, "invalid auth header")
			return
		}
		userId, err := h.services.User.ParseToken(headerParts)
		if err != nil {
			newErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		blackUser := h.blackList.IsUserBlackListed(userId)
		if blackUser == true {
			newErrorResponse(w, http.StatusUnauthorized, "invalid auth header")
			return
		}

		ctx := context.WithValue(r.Context(), userCtx, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
