package model

import "github.com/jinzhu/gorm"

type Organization struct {
	gorm.Model
	Name string `gorm:"not null"`
}

type OrganizationsAgreements struct {
	gorm.Model
	FromID           uint
	ToID             uint
	OrganizationFrom Organization
	OrganizationTo   Organization
}
