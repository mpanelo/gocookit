package middleware

import (
	"net/http"
	"strings"

	"github.com/mpanelo/gocookit/context"
	"github.com/mpanelo/gocookit/models"
)

type User struct {
	models.UserService
}

func (u *User) Apply(next http.Handler) http.HandlerFunc {
	return u.ApplyFn(next.ServeHTTP)
}

func (u *User) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.HasPrefix(path, "/assets/") || strings.HasPrefix(path, "/images/") {
			next(rw, r)
			return
		}

		cookie, err := r.Cookie("remember_token")
		if err != nil {
			next(rw, r)
			return
		}
		user, err := u.ByRemember(cookie.Value)
		if err != nil {
			next(rw, r)
			return
		}

		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)

		next(rw, r)
	}
}

type RequireUser struct {
}

func (ru *RequireUser) Apply(next http.Handler) http.HandlerFunc {
	return ru.ApplyFn(next.ServeHTTP)
}

func (ru *RequireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			http.Redirect(rw, r, "/signin", http.StatusFound)
			return
		}

		next(rw, r)
	}
}
