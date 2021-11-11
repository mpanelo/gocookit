package controllers

import (
	"fmt"
	"net/http"

	"github.com/mpanelo/gocookit/models"
	"github.com/mpanelo/gocookit/views"
)

func NewUsers(us models.UserService) *Users {
	return &Users{
		SignUp: views.NewView("users/signup"),
		us:     us,
	}
}

type Users struct {
	SignUp *views.View
	us     models.UserService
}

type SignUpForm struct {
	Name     string
	Email    string
	Password string
}

func (u *Users) Create(rw http.ResponseWriter, r *http.Request) {
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

	user, err = u.us.ByEmail(form.Email)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError) // TODO display user friendly error on page
		return
	}
	fmt.Fprintln(rw, user)
}
