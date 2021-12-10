package models

import (
	"regexp"
	"strings"

	"github.com/mpanelo/gocookit/hash"
	"github.com/mpanelo/gocookit/rand"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name         string `gorm:"not null"`
	Email        string `gorm:"not null;uniqueIndex"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null"`
}

type UserService interface {
	UserDB
	Authenticate(string, string) (*User, error)
}

type userService struct {
	UserDB
	pepper string
}

func NewUserService(db *gorm.DB, hmacKey, pepper string) UserService {
	return &userService{
		UserDB: &userValidator{
			UserDB:     &userGorm{db},
			hmac:       hash.NewHmac(hmacKey),
			emailRegex: regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"),
		},
		pepper: pepper,
	}
}

func (us *userService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.UserDB.ByEmail(email)
	if err != nil {
		if err == ErrNotFound {
			return nil, ErrUserCredentialsInvalid
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password+us.pepper))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, ErrUserCredentialsInvalid
		}
		return nil, err
	}

	return foundUser, nil
}

type UserDB interface {
	ByID(uint) (*User, error)
	ByEmail(string) (*User, error)
	ByRemember(string) (*User, error)
	Create(*User) error
	Update(*User) error
}

type userValidator struct {
	UserDB
	hmac       *hash.Hmac
	emailRegex *regexp.Regexp
	pepper     string
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

func (uv *userValidator) ByRemember(remember string) (*User, error) {
	var user User
	user.Remember = remember

	err := runUserValidatorFuncs(&user,
		uv.rememberTokenMinLength,
		uv.generateRememberHash,
		uv.requireRememberHash)
	if err != nil {
		return nil, err
	}

	return uv.UserDB.ByRemember(user.RememberHash)
}

func (uv *userValidator) Create(user *User) error {
	err := runUserValidatorFuncs(user,
		uv.requirePassword,
		uv.passwordMinLength,
		uv.generatePasswordHash,
		uv.requirePasswordHash,
		uv.setRememberIfUnset,
		uv.rememberTokenMinLength,
		uv.generateRememberHash,
		uv.requireRememberHash,
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

func (uv *userValidator) Update(user *User) error {
	err := runUserValidatorFuncs(user,
		uv.passwordMinLength,
		uv.generatePasswordHash,
		uv.requirePasswordHash,
		uv.rememberTokenMinLength,
		uv.generateRememberHash,
		uv.requireRememberHash,
		uv.normalizeEmail,
		uv.validateEmailFormat,
		uv.emailIsAvail)
	if err != nil {
		return err
	}

	return uv.UserDB.Update(user)
}

type userValidatorFunc func(*User) error

func (uv *userValidator) requirePassword(user *User) error {
	if user.Password == "" {
		return ErrUserPasswordRequired
	}
	return nil
}

func (uv *userValidator) passwordMinLength(user *User) error {
	if user.Password == "" {
		return nil
	}
	if len(user.Password) < 8 {
		return ErrUserPasswordTooShort
	}
	return nil
}

func (uv *userValidator) generatePasswordHash(user *User) error {
	if user.Password == "" {
		return nil
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password+uv.pepper), bcrypt.DefaultCost)
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

func (uv *userValidator) setRememberIfUnset(user *User) error {
	if user.Remember != "" {
		return nil
	}

	token, err := rand.RememberToken()
	if err != nil {
		return err
	}

	user.Remember = token
	return nil
}

func (uv *userValidator) rememberTokenMinLength(user *User) error {
	if user.Remember == "" {
		return nil
	}

	n, err := rand.NBytes(user.Remember)
	if err != nil {
		return err
	}

	if n < rand.RememberTokenBytesLen {
		return ErrUserRememberTooShort
	}
	return nil
}

func (uv *userValidator) generateRememberHash(user *User) error {
	if user.Remember != "" {
		user.RememberHash = uv.hmac.Hash(user.Remember)
	}
	return nil
}

func (uv *userValidator) requireRememberHash(user *User) error {
	if user.RememberHash == "" {
		return ErrUserRememberHashRequired
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
	if !uv.emailRegex.MatchString(user.Email) {
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

func (ug *userGorm) ByRemember(rememberHash string) (*User, error) {
	var user User
	tx := ug.db.Where("remember_hash", rememberHash)

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

func (ug *userGorm) Update(user *User) error {
	result := ug.db.Save(user)
	return result.Error
}
