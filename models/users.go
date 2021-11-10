package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Email    string `gorm:"not null;uniqueIndex"`
	Password string `gorm:"-"`
}

type UserService interface {
	UserDB
}

type userService struct {
	UserDB
}

func NewUserService(db *gorm.DB) UserService {
	return &userService{&userValidator{&userGorm{db}}}
}

type UserDB interface {
	Create(*User) error
}

type userValidator struct {
	UserDB
}

type userGorm struct {
	db *gorm.DB
}

func (ug *userGorm) Create(user *User) error {
	result := ug.db.Create(user)
	return result.Error
}
