package main

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/BlockTeam4Boys/digitaldocs/internal/token"
)

func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(JWTInfo.RefreshToken.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	tokenString := cookie.Value
	claims, err := token.RefreshTokenInfo(tokenString, SessionRepo.GetUserKey)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	sessionID := claims.SessionID
	dbRefreshTokenHash, err := SessionRepo.GetRefreshTokenHash(sessionID)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	refreshTokenHash := hash(tokenString)
	userID := claims.UserID
	key, err := SessionRepo.GetUserKey(userID)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if refreshTokenHash == dbRefreshTokenHash {
		refreshExpirationTime := time.Now().Add(JWTInfo.RefreshToken.TTL)
		accessExpirationTime := time.Now().Add(JWTInfo.AccessToken.TTL)
		refreshClaims := token.RefreshTokenClaims{
			UserID:                 userID,
			SessionID:              sessionID,
			ParentRefreshTokenHash: refreshTokenHash,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: refreshExpirationTime.Unix(),
			},
		}
		refreshToken, err := token.New(refreshClaims, []byte(key))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		curRefreshTokenHash := hash(refreshToken)
		accessClaims := token.AccessTokenClaims{
			UserID:                  userID,
			SessionID:               sessionID,
			ParentRefreshTokenHash:  refreshTokenHash,
			CurrentRefreshTokenHash: curRefreshTokenHash,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: accessExpirationTime.Unix(),
			},
		}
		accessToken, err := token.New(accessClaims, []byte(key))
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
		return
	}
	if claims.ParentRefreshTokenHash != dbRefreshTokenHash {
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
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	SessionRepo.SetRefreshTokenHash(sessionID, refreshTokenHash)
	refreshExpirationTime := time.Now().Add(JWTInfo.RefreshToken.TTL)
	accessExpirationTime := time.Now().Add(JWTInfo.AccessToken.TTL)
	refreshClaims := token.RefreshTokenClaims{
		UserID:                 userID,
		SessionID:              sessionID,
		ParentRefreshTokenHash: refreshTokenHash,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshExpirationTime.Unix(),
		},
	}
	refreshToken, err := token.New(refreshClaims, []byte(key))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	curRefreshTokenHash := hash(refreshToken)
	accessClaims := token.AccessTokenClaims{
		UserID:                  userID,
		SessionID:               sessionID,
		ParentRefreshTokenHash:  refreshTokenHash,
		CurrentRefreshTokenHash: curRefreshTokenHash,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: accessExpirationTime.Unix(),
		},
	}
	accessToken, err := token.New(accessClaims, []byte(key))
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
	return
}
