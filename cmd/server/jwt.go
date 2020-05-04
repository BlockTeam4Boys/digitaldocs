package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	TokenCookieName  string
	TokenTTL         time.Duration
	TokenRefreshTime time.Duration
	JWTKey           []byte

	TokenExpiredErr = errors.New("token expired")
	InvalidTokenErr = errors.New("invalid token")
)

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

func CheckJWT(r *http.Request) (Claims, error) {
	cookie, err := r.Cookie(TokenCookieName)
	claims := Claims{}
	if err != nil {
		return claims, err
	}
	tokenString := cookie.Value
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return JWTKey, nil
	})
	if err != nil {
		return claims, err
	}
	if !token.Valid {
		return claims, InvalidTokenErr
	}
	if time.Now().Sub(time.Unix(claims.ExpiresAt, 0)) > TokenRefreshTime {
		return claims, TokenExpiredErr
	}
	return claims, nil
}
