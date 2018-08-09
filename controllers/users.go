package controllers

import (
	"fmt"
	"net/http"

	"github.com/tjcain/theFieldBiologist/views"
)

// Users represents a new user view
type Users struct {
	NewView *views.View
}

// NewUsers is used to create a new Users controller.
// Any error in rendering templates will cause this function to panic.
func NewUsers() *Users {
	return &Users{
		NewView: views.NewView("pages", "users/signup"),
	}
}

// New renders a new user view.
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}

// Create is a handler function responsible for handing the post request
// made by our signup form
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var form SignUpForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	fmt.Fprintln(w, form)
}

// SignUpForm stores data POSTed from our signup form
type SignUpForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}
