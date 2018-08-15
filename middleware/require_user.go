package middleware

import (
	"net/http"
	"strings"

	"github.com/tjcain/theFieldBiologist/context"
	"github.com/tjcain/theFieldBiologist/models"
)

// User middleware looks up the current user by remember_token. If the user is
// found, they will be set on the request context.
// Regardless of result, the next handler is always called.
type User struct {
	models.UserService
}

// Apply ...
func (mw *User) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

// ApplyFn ...
func (mw *User) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.HasPrefix(path, "/assets/") {
			next(w, r)
			return
		}
		cookie, err := r.Cookie("remember_token")
		if err != nil {
			next(w, r)
			return
		}
		user, err := mw.UserService.ByRemember(cookie.Value)
		if err != nil {
			next(w, r)
			return
		}
		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)
		next(w, r)
	})
}

// RequireUser contains the UserService
type RequireUser struct {
	models.UserService
}

// Apply ...
func (mw *RequireUser) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

// ApplyFn will return a http.HandlerFunc that will check to see if a user
// is logged in and then eitehr call next(w, r) if they are, or redirect them to
// the login page if they are not.
func (mw *RequireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next(w, r)
	})
}
