package models

import (
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	sharedSecretPepper = "pepper" // TODO load from an environment variable
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type User struct {
	gorm.Model
	Name         string `gorm:"not null"`
	Email        string `gorm:"not null;uniqueIndex"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
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
	ByID(uint) (*User, error)
	ByEmail(string) (*User, error)
	Create(*User) error
}

type userValidator struct {
	UserDB
}

func (uv *userValidator) ByEmail(email string) (*User, error) {
	var user User
	user.Email = email

	err := runUserValidatorFuncs(&user, normalizeEmail)
	if err != nil {
		return nil, err
	}

	return uv.UserDB.ByEmail(user.Email)
}

func (uv *userValidator) Create(user *User) error {
	err := runUserValidatorFuncs(user,
		requirePassword,
		passwordMinLength,
		generatePasswordHash,
		requirePasswordHash,
		requireEmail,
		normalizeEmail,
		validateEmailFormat,
		requireName)
	if err != nil {
		return err
	}

	return uv.UserDB.Create(user)
}

type userValidatorFunc func(*User) error

func requirePassword(user *User) error {
	if user.Password == "" {
		return ErrUserPasswordRequired
	}
	return nil
}

func passwordMinLength(user *User) error {
	if len(user.Password) < 8 {
		return ErrUserPasswordTooShort
	}
	return nil
}

func generatePasswordHash(user *User) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password+sharedSecretPepper), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(passwordHash)
	user.Password = ""
	return nil
}

func requirePasswordHash(user *User) error {
	if user.PasswordHash == "" {
		return ErrUserPasswordHashRequired
	}
	return nil
}

func requireEmail(user *User) error {
	if user.Email == "" {
		return ErrUserEmailRequired
	}
	return nil
}

func normalizeEmail(user *User) error {
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)
	return nil
}

func validateEmailFormat(user *User) error {
	if !emailRegex.MatchString(user.Email) {
		return ErrUserEmailInvalid
	}
	return nil
}

func requireName(user *User) error {
	if user.Name == "" {
		return ErrUserNameRequired
	}
	return nil
}

func runUserValidatorFuncs(user *User, funcs ...userValidatorFunc) error {
	for _, f := range funcs {
		if err := f(user); err != nil {
			return err
		}
	}

	return nil
}

type userGorm struct {
	db *gorm.DB
}

func (ug *userGorm) ByID(id uint) (*User, error) {
	return first(ug.db.Where("id = ?", id))
}

func (ug *userGorm) ByEmail(email string) (*User, error) {
	return first(ug.db.Where("email = ?", email))
}

func first(tx *gorm.DB) (*User, error) {
	var user User

	err := tx.First(&user).Error
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (ug *userGorm) Create(user *User) error {
	result := ug.db.Create(user)
	return result.Error
}
