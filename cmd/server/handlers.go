package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"

	"github.com/BlockTeam4Boys/digitaldocs/internal/model"
)

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	//Extract fields
	r.ParseForm()
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")
	if password == "" || email == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//Find in DB
	passwordHash := sha256.Sum256([]byte(password))
	var user model.User
	if DB.Where(&model.User{
		Email:        email,
		PasswordHash: hex.EncodeToString(passwordHash[:]),
	}).First(&user).RecordNotFound() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	//Generate token
	expirationTime := time.Now().Add(TokenTTL)
	claims := &Claims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JWTKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     TokenCookieName,
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true,
	})
}

func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	//Check token
	claims, err := CheckJWT(r)
	switch err {
	case nil:
		break
	case http.ErrNoCookie, jwt.ErrSignatureInvalid, InvalidTokenErr:
		w.WriteHeader(http.StatusUnauthorized)
		return
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//Renew
	expirationTime := time.Now().Add(TokenTTL)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	tokenString, err := token.SignedString(JWTKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    TokenCookieName,
		Value:   tokenString,
		Expires: expirationTime,
	})
}

func RequestHandler(w http.ResponseWriter, r *http.Request) {
	//Check token
	claims, err := CheckJWT(r)
	switch err {
	case nil:
		break
	case http.ErrNoCookie, jwt.ErrSignatureInvalid, InvalidTokenErr:
		w.WriteHeader(http.StatusUnauthorized)
		return
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//Parse values
	r.ParseForm()
	organizationIDAsStr := r.PostForm.Get("organization")
	diplomaIDAsStr := r.PostForm.Get("diploma")
	organizationID, err := strconv.Atoi(organizationIDAsStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	diplomaID, err := strconv.Atoi(diplomaIDAsStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//Find agreement
	var agreement model.OrganizationsAgreements
	notFound := DB.Where("from_id = ? AND to_id = ?",
		uint(organizationID),
		DB.Table("users").
			Select("organization_id").
			Where(&model.User{
				Model: gorm.Model{
					ID: claims.UserID,
				},
			}).SubQuery()).First(&agreement).RecordNotFound()
	if notFound {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//Request diploma
	diploma := fmt.Sprintf(`cool diploma №%v`, diplomaID) //TODO: make diploma request
	//Send response
	w.Write([]byte(diploma))
}
