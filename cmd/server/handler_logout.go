package main

import (
	"net/http"

	"github.com/BlockTeam4Boys/digitaldocs/internal/middleware"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := r.Context().Value(middleware.SessionIDKey).(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	SessionRepo.DeleteRefreshToken(sessionID)
	http.SetCookie(w, &http.Cookie{
		Name:   JWTInfo.AccessToken.CookieName,
		Value:  "",
		MaxAge: -1,
	})
	http.SetCookie(w, &http.Cookie{
		Name:   JWTInfo.RefreshToken.CookieName,
		Value:  "",
		MaxAge: -1,
	})
}
