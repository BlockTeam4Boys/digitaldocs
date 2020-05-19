package repository

import (
	"github.com/jinzhu/gorm"

	"github.com/BlockTeam4Boys/digitaldocs/internal/model"
)

type PostgresUserRepository struct {
	db *gorm.DB
}

func NewPostgresUserRepository(db *gorm.DB) *PostgresUserRepository {
	return &PostgresUserRepository{
		db: db,
	}
}

func (p *PostgresUserRepository) Find(email string) (model.User, error) {
	var user model.User
	err := p.db.Where(&model.User{
		Email: email,
	}).First(&user).Error
	return user, err
}
