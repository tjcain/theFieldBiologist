package middleware

import (
	"net/http"

	"github.com/tjcain/theFieldBiologist/context"
	"github.com/tjcain/theFieldBiologist/models"
)

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
		cookie, err := r.Cookie("remember_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		user, err := mw.UserService.ByRemember(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)
		next(w, r)
	})
}
