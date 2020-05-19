package token

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

var (
	InvalidTokenErr = errors.New("invalid token")
)

type AccessTokenClaims struct {
	SessionID               string
	UserID                  uint
	ParentRefreshTokenHash  string
	CurrentRefreshTokenHash string
	jwt.StandardClaims
}

type RefreshTokenClaims struct {
	UserID                 uint
	SessionID              string
	ParentRefreshTokenHash string
	jwt.StandardClaims
}

type getKeyFunc func(uint) (string, error)

func AccessTokenInfo(tokenStr string, getKey getKeyFunc) (AccessTokenClaims, error) {
	claims := AccessTokenClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		userID := claims.UserID
		key, err := getKey(userID)
		return []byte(key), err
	})
	if err != nil {
		return claims, err
	}
	if !token.Valid {
		return claims, InvalidTokenErr
	}
	return claims, nil
}

func RefreshTokenInfo(tokenString string, getKey getKeyFunc) (RefreshTokenClaims, error) {
	claims := RefreshTokenClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		userID := claims.UserID
		key, err := getKey(userID)
		return []byte(key), err
	})
	if err != nil {
		return claims, err
	}
	if !token.Valid {
		return claims, InvalidTokenErr
	}
	return claims, nil
}

func New(claims jwt.Claims, key []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(key)
	return tokenString, err
}
