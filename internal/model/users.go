package model

import "github.com/jinzhu/gorm"

type Role struct {
	ID   uint   `gorm:"primary_key"`
	Name string `gorm:"unique;not null"`
}

type User struct {
	gorm.Model
	Email          string `gorm:"unique;not null"`
	PasswordHash   string `gorm:"not null"`
	OrganizationID uint
	Organization   Organization
	RoleID         uint
	Role           Role
}
