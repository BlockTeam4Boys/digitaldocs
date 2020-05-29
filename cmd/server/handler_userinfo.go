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
	const query = `
		SELECT u.email, 
			u.organization_id,
			r.Name as role,
			o1.name as organization_name, 
			oa1.to_id, 
			o2.name as to_name, 
			oa2.from_id, 
			o3.name as from_name 
		FROM   users AS u
			join roles AS r
			ON u.role_id = r.id
			join organizations AS o1 
			ON u.organization_id = o1.id 
			join organizations_agreements AS oa1 
			ON u.organization_id = oa1.from_id 
			join organizations AS o2 
			ON oa1.to_id = o2.id 
			join organizations_agreements AS oa2 
			ON u.organization_id = oa2.to_id 
			join organizations AS o3 
			ON oa2.from_id = o3.id 
		WHERE  u.id = ?;`
	var result []struct {
		Email            string
		Role             string
		OrganizationID   uint
		OrganizationName string
		ToID             uint
		ToName           string
		FromID           uint
		FromName         string
	}
	DB.Raw(query, userID).Scan(&result)
	var info UserInfo
	if len(result) == 0 {
		return info
	}
	info.Email = result[0].Email
	info.OrganizationID = result[0].OrganizationID
	info.OrganizationName = result[0].OrganizationName
	info.Role = result[0].Role
	info.AllowAccessFrom = make(map[uint]string)
	info.AllowAccessTo = make(map[uint]string)
	for i := range result {
		info.AllowAccessTo[result[i].ToID] = result[i].ToName
		info.AllowAccessFrom[result[i].FromID] = result[i].FromName
	}
	return info
}
