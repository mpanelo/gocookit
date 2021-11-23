package controllers

import (
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
	var vd views.Data
	var form SignUpForm

	if err := parseForm(r, &form); err != nil {
		vd.SetAlertDanger(err)
		u.SignUpView.Render(rw, r, vd)
		return
	}

	user := &models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}

	err := u.us.Create(user)
	if err != nil {
		vd.SetAlertDanger(err)
		u.SignUpView.Render(rw, r, vd)
		return
	}

	err = u.setRememberTokenCookie(rw, user)
	if err != nil {
		vd.SetAlertDanger(err)
		u.SignInView.Render(rw, r, vd)
		return
	}

	http.Redirect(rw, r, "/recipes", http.StatusFound)
}

type SignInForm struct {
	Email    string
	Password string
}

func (u *Users) SignIn(rw http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form SignInForm

	if err := parseForm(r, &form); err != nil {
		vd.SetAlertDanger(err)
		u.SignInView.Render(rw, r, vd)
		return
	}

	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		vd.SetAlertDanger(err)
		u.SignInView.Render(rw, r, vd)
		return
	}

	err = u.setRememberTokenCookie(rw, user)
	if err != nil {
		vd.SetAlertDanger(err)
		u.SignInView.Render(rw, r, vd)
		return
	}

	http.Redirect(rw, r, "/recipes", http.StatusFound)
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
