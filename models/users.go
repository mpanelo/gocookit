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

	err := runUserValidatorFuncs(&user, uv.normalizeEmail)
	if err != nil {
		return nil, err
	}

	return uv.UserDB.ByEmail(user.Email)
}

func (uv *userValidator) Create(user *User) error {
	err := runUserValidatorFuncs(user,
		uv.requirePassword,
		uv.passwordMinLength,
		uv.generatePasswordHash,
		uv.requirePasswordHash,
		uv.requireEmail,
		uv.normalizeEmail,
		uv.validateEmailFormat,
		uv.emailIsAvail,
		uv.requireName)
	if err != nil {
		return err
	}

	return uv.UserDB.Create(user)
}

type userValidatorFunc func(*User) error

func (uv *userValidator) requirePassword(user *User) error {
	if user.Password == "" {
		return ErrUserPasswordRequired
	}
	return nil
}

func (uv *userValidator) passwordMinLength(user *User) error {
	if len(user.Password) < 8 {
		return ErrUserPasswordTooShort
	}
	return nil
}

func (uv *userValidator) generatePasswordHash(user *User) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password+sharedSecretPepper), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(passwordHash)
	user.Password = ""
	return nil
}

func (uv *userValidator) requirePasswordHash(user *User) error {
	if user.PasswordHash == "" {
		return ErrUserPasswordHashRequired
	}
	return nil
}

func (uv *userValidator) requireEmail(user *User) error {
	if user.Email == "" {
		return ErrUserEmailRequired
	}
	return nil
}

func (uv *userValidator) normalizeEmail(user *User) error {
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)
	return nil
}

func (uv *userValidator) validateEmailFormat(user *User) error {
	if !emailRegex.MatchString(user.Email) {
		return ErrUserEmailInvalid
	}
	return nil
}

func (uv *userValidator) emailIsAvail(user *User) error {
	foundUser, err := uv.UserDB.ByEmail(user.Email)
	if err != nil {
		if err == ErrNotFound {
			return nil
		}
		return err
	}

	if foundUser.ID != user.ID {
		return ErrUserEmailTaken
	}

	return nil
}

func (uv *userValidator) requireName(user *User) error {
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
	var user User
	tx := ug.db.Where("id = ?", id)

	err := first(tx, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	tx := ug.db.Where("email = ?", email)

	err := first(tx, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func first(tx *gorm.DB, dst interface{}) error {
	err := tx.First(dst).Error
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}

func (ug *userGorm) Create(user *User) error {
	result := ug.db.Create(user)
	return result.Error
}
