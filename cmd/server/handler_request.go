package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"

	"github.com/BlockTeam4Boys/digitaldocs/internal/middleware"
	"github.com/BlockTeam4Boys/digitaldocs/internal/model"
)

func RequestHandler(w http.ResponseWriter, r *http.Request) {
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
	//Find userID in context
	userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
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
					ID: userID,
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
