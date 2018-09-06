package middleware

import (
	"net/http"

	"github.com/tjcain/theFieldBiologist/context"
	"github.com/tjcain/theFieldBiologist/models"
)

// RequireAdmin contains the UserService
type RequireAdmin struct {
	models.UserService
}

// Apply ...
func (mw *RequireAdmin) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

// ApplyFn will return a http.HandlerFunc that will check to see if a user
// has the admin role in and then either call next(w, r) if they are, or
// redirect them to the home page if they are not.
func (mw *RequireAdmin) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if !user.Admin {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		next(w, r)
	})
}
