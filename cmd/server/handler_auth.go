package main

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	guuid "github.com/google/uuid"

	"github.com/BlockTeam4Boys/digitaldocs/internal/token"
)

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")
	if password == "" || email == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	passwordHash := hash(password)
	user, err := UserRepo.Find(email)
	if err != nil || user.PasswordHash != passwordHash {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	//Generate tokens
	key := []byte(passwordHash)
	sessionID := guuid.New().String()
	refreshExpirationTime := time.Now().Add(JWTInfo.RefreshToken.TTL)
	refreshClaims := token.RefreshTokenClaims{
		UserID:                 user.ID,
		SessionID:              sessionID,
		ParentRefreshTokenHash: "",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshExpirationTime.Unix(),
		},
	}
	refreshToken, err := token.New(refreshClaims, key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	accessExpirationTime := time.Now().Add(JWTInfo.AccessToken.TTL)
	refreshTokenHash := hash(refreshToken)
	accessClaims := token.AccessTokenClaims{
		UserID:                  user.ID,
		SessionID:               sessionID,
		ParentRefreshTokenHash:  "",
		CurrentRefreshTokenHash: refreshTokenHash,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: accessExpirationTime.Unix(),
		},
	}
	accessToken, err := token.New(accessClaims, key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//Save values
	SessionRepo.SetRefreshTokenHash(sessionID, refreshTokenHash)
	SessionRepo.SetUserKey(user.ID, key)
	//Set cookie
	http.SetCookie(w, &http.Cookie{
		Name:     JWTInfo.AccessToken.CookieName,
		Value:    accessToken,
		Expires:  accessExpirationTime,
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     JWTInfo.RefreshToken.CookieName,
		Value:    refreshToken,
		Expires:  refreshExpirationTime,
		HttpOnly: true,
	})
}
