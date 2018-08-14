package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/tjcain/theFieldBiologist/models"
	"github.com/tjcain/theFieldBiologist/rand"
	"github.com/tjcain/theFieldBiologist/views"
)

// SignUpForm stores data POSTed from our signup form
type SignUpForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// LogInForm stores data POSTed from our login form
type LogInForm struct {
	Email      string `schema:"email"`
	Password   string `schema:"password"`
	RememberMe bool   `schema:"remember_me"`
}

// Users represents a new user view
type Users struct {
	NewView         *views.View
	LogInView       *views.View
	AllArticlesView *views.View
	us              models.UserService
}

// NewUsers is used to create a new Users controller.
// Any error in rendering templates will cause this function to panic.
func NewUsers(us models.UserService) *Users {
	return &Users{
		NewView:         views.NewView("pages", "users/signup"),
		LogInView:       views.NewView("pages", "users/login"),
		AllArticlesView: views.NewView("pages", "users/allarticles"),
		us:              us,
	}
}

// New renders a new user view.
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	u.NewView.Render(w, nil)
}

// Create is a handler function responsible for handing the post request
// made by our signup form
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var form SignUpForm
	var vd views.Data
	if err := parseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		u.NewView.Render(w, vd)
		return
	}
	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}
	if err := u.us.Create(&user); err != nil {
		vd.SetAlert(err)
		u.NewView.Render(w, vd)
		return
	}
	err := u.signIn(w, &user)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	// Redirect to the cookie test page to test the cookie
	http.Redirect(w, r, "/articles", http.StatusFound)
}

// signIn is used to sign the given user in via cookies
func (u *Users) signIn(w http.ResponseWriter, user *models.User) error {
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
	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    user.Remember,
		HttpOnly: true,
	}
	if user.RememberMe {
		cookie.Expires = time.Now().Add(time.Hour * 336)
	}
	http.SetCookie(w, &cookie)
	return nil
}

// LogIn processes a login form when a user attempts to log in with a
// email and password
func (u *Users) LogIn(w http.ResponseWriter, r *http.Request) {
	var form LogInForm
	var vd views.Data
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.LogInView.Render(w, vd)
		return
	}
	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			vd.AlertError("No user exists with that email address")
		default:
			vd.SetAlert(err)
		}
		u.LogInView.Render(w, vd)
		return
	}
	user.RememberMe = form.RememberMe
	err = u.signIn(w, user)
	if err != nil {
		vd.SetAlert(err)
		u.LogInView.Render(w, vd)
		return
	}
	http.Redirect(w, r, "/articles", http.StatusFound)
}

// ShowAllArticles lists all the articles belonging to a given author, this is
// only visible to a logged in author as it provides edit / delete capability
func (u *Users) ShowAllArticles(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("remember_token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user, err := u.us.ByRemember(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var vd views.Data
	fmt.Println(user)
	articles, err := u.us.ArticlesByUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	vd.User = user
	vd.Yield = articles
	u.AllArticlesView.Render(w, vd)

}

// CookieTest is a temporary function for development only. It will display
// the cookies set on a current user.
func (u *Users) CookieTest(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("remember_token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user, err := u.us.ByRemember(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, user)
}
