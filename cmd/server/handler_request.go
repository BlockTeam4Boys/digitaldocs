package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"

	"github.com/BlockTeam4Boys/digitaldocs/internal/middleware"
	"github.com/BlockTeam4Boys/digitaldocs/internal/model"
)

func DocumentHandler(w http.ResponseWriter, r *http.Request) {
	//Parse values
	r.ParseForm()
	organizationIDAsStr := r.PostForm.Get("universityId")
	number := r.PostForm.Get("number")
	firstName := r.PostForm.Get("firstName")
	secondName := r.PostForm.Get("secondName")
	spec := r.PostForm.Get("spec")
	year := r.PostForm.Get("year")

	organizationID, err := strconv.Atoi(organizationIDAsStr)
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
	diplomaRequestBody := map[string]interface{}{
		"number":     number,
		"firstName":  firstName,
		"secondName": secondName,
		"spec":       spec,
		"year":       year,
	}

	bytesRepresentation, err := json.Marshal(diplomaRequestBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var url = "http://127.0.0.1:"
	if organizationID == 2 {
		url = url + "8888"
	} else {
		url = url + "8887"
	}
	url = url + "/diploma"

	diplomaResponse, err := http.Post(url, "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	diploma, err := ioutil.ReadAll(diplomaResponse.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write(diploma)
}
