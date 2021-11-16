package controllers

import (
	"fmt"
	"net/http"

	"github.com/mpanelo/gocookit/models"
	"github.com/mpanelo/gocookit/rand"
	"github.com/mpanelo/gocookit/views"
)

func NewUsers(us models.UserService) *Users {
	return &Users{
		SignUpView: views.NewView("users/signup"),
		SignInView: views.NewView("users/signin"),
		us:         us,
	}
}

type Users struct {
	SignUpView *views.View
	SignInView *views.View
	us         models.UserService
}

type SignUpForm struct {
	Name     string
	Email    string
	Password string
}

func (u *Users) SignUp(rw http.ResponseWriter, r *http.Request) {
	var form SignUpForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}

	user := &models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}

	err := u.us.Create(user)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError) // TODO display user friendly error on page
		return
	}

	err = u.setRememberTokenCookie(rw, user)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError) // TODO display user friendly error on page
		return
	}

	http.Redirect(rw, r, "/whoami", http.StatusFound)
}

type SignInForm struct {
	Email    string
	Password string
}

func (u *Users) SignIn(rw http.ResponseWriter, r *http.Request) {
	var form SignInForm
	if err := parseForm(r, &form); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError) // TODO display user friendly error on page
		return
	}

	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError) // TODO display user friendly error on page
		return
	}

	err = u.setRememberTokenCookie(rw, user)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError) // TODO display user friendly error on page
		return
	}

	http.Redirect(rw, r, "/whoami", http.StatusFound)
}

func (u *Users) setRememberTokenCookie(rw http.ResponseWriter, user *models.User) error {
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}

		user.Remember = token

		err = u.us.Update(user)
		if err != nil {
			return err
		}
	}

	http.SetCookie(rw, &http.Cookie{
		Name:     "remember_token",
		Value:    user.Remember,
		HttpOnly: true,
	})

	return nil
}

func (u *Users) Whoami(rw http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("remember_token")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := u.us.ByRemember(cookie.Value)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(rw, "You are "+user.Name)
}
