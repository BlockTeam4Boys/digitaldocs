package main

import (
	"encoding/json"
	"net/http"

	"github.com/BlockTeam4Boys/digitaldocs/internal/middleware"
)

func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	userInfo := getUserInfo(userID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userInfo)
}

type UserInfo struct {
	Email            string
	Role             string
	OrganizationID   uint
	OrganizationName string
	AllowAccessFrom  map[uint]string
	AllowAccessTo    map[uint]string
}

func getUserInfo(userID uint) UserInfo {
	var user struct {
		Email            string
		Role             string
		OrganizationID   uint
		OrganizationName string
	}

	const userQuery = `
		SELECT u.id,
       		u.email,
       		u.organization_id,
       		r.Name as role,
       		o.name as organization_name
		FROM users AS u
			join roles AS r
				ON u.role_id = r.id
         	join organizations AS o
				ON u.organization_id = o.id
		WHERE u.id = ?;`

	DB.Raw(userQuery, userID).Scan(&user)
	var info UserInfo
	//todo: check user for nil
	info.Email = user.Email
	info.OrganizationID = user.OrganizationID
	info.OrganizationName = user.OrganizationName
	info.Role = user.Role

	var organization []struct {
		ID   uint
		Name string
	}

	const allowAccessToQuery = `
		SELECT o.id, 
				o.name
		FROM organizations_agreements AS oa
			left join organizations o 
				ON oa.to_id = o.id
		WHERE oa.from_id = ?;`

	DB.Raw(allowAccessToQuery, user.OrganizationID).Scan(&organization)
	info.AllowAccessTo = make(map[uint]string)
	for i := range organization {
		info.AllowAccessTo[organization[i].ID] = organization[i].Name
	}

	const allowAccessFromQuery = `
		SELECT o.id, 
				o.name
		FROM organizations_agreements AS oa
			left join organizations o 
				ON oa.from_id = o.id
		WHERE oa.to_id = ?;`

	DB.Raw(allowAccessFromQuery, user.OrganizationID).Scan(&organization)
	info.AllowAccessFrom = make(map[uint]string)
	for i := range organization {
		info.AllowAccessFrom[organization[i].ID] = organization[i].Name
	}

	return info
}
