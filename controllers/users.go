package controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/tjcain/theFieldBiologist/context"
	"github.com/tjcain/theFieldBiologist/email"
	"github.com/tjcain/theFieldBiologist/models"
	"github.com/tjcain/theFieldBiologist/rand"
	"github.com/tjcain/theFieldBiologist/views"
)

const (
	// ManageArticles is the named route for /user/articles
	ManageArticles = "mange_articles"
)

// SignUpForm stores data POSTed from our signup form
type SignUpForm struct {
	Name      string `schema:"name"`
	Email     string `schema:"email"`
	Password  string `schema:"password"`
	TandC     bool   `schema:"tandc"`
	EmailPerm bool   `schema:"emailpermission"`
}

// LogInForm stores data POSTed from our login form
type LogInForm struct {
	Email      string `schema:"email"`
	Password   string `schema:"password"`
	RememberMe bool   `schema:"remember_me"`
}

// EditProfileForm stores data POSTed from our edit profile form
type EditProfileForm struct {
	Email string `schema:"email"`
	Name  string `schema:"name"`
	Bio   string `schema:"bio"`
}

// ResetPwForm stores data POSTed during the reset password process
type ResetPwForm struct {
	Email    string `schema:"email"`
	Token    string `schema:"token"`
	Password string `schema:"password"`
}

// Users represents a new user view
type Users struct {
	NewView         *views.View
	LogInView       *views.View
	AllArticlesView *views.View
	UserProfileView *views.View
	ProfileView     *views.View
	ForgotPwView    *views.View
	ResetPwView     *views.View
	us              models.UserService
	emailer         *email.Client
}

// NewUsers is used to create a new Users controller.
// Any error in rendering templates will cause this function to panic.
func NewUsers(us models.UserService, emailer *email.Client) *Users {
	return &Users{
		NewView:         views.NewView("pages", "users/signup"),
		LogInView:       views.NewView("pages", "users/login"),
		AllArticlesView: views.NewView("pages", "users/allarticles"),
		UserProfileView: views.NewView("pages", "users/editprofile"),
		ProfileView:     views.NewView("pages", "users/profileview"),
		ForgotPwView:    views.NewView("pages", "users/forgot_pw"),
		ResetPwView:     views.NewView("pages", "users/reset_pw"),
		us:              us,
		emailer:         emailer,
	}
}

// New renders a new user view.
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	u.NewView.Render(w, r, nil)
}

// Create is a handler function responsible for handing the post request
// made by our signup form
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form SignUpForm
	vd.Yield = &form
	// fmt.Printf("%+v\n", vd.Yield)
	if err := parseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		u.NewView.Render(w, r, vd)
		return
	}
	user := models.User{
		Name:            form.Name,
		Email:           form.Email,
		Password:        form.Password,
		TandC:           form.TandC,
		EmailPermission: form.EmailPerm,
	}
	if err := u.us.Create(&user); err != nil {
		vd.SetAlert(err)
		u.NewView.Render(w, r, vd)
		return
	}
	err := u.signIn(w, &user)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	// Redirect to the cookie test page to test the cookie
	http.Redirect(w, r, "/user/articles", http.StatusFound)
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

// LogIn POST /login
func (u *Users) LogIn(w http.ResponseWriter, r *http.Request) {
	var form LogInForm
	var vd views.Data
	vd.Yield = &form

	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.LogInView.Render(w, r, vd)
		return
	}
	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			fmt.Printf("%+v\n", vd.Yield)
			vd.AlertError("No user exists with that email address")
		default:
			vd.SetAlert(err)
		}
		u.LogInView.Render(w, r, vd)
		return
	}
	user.RememberMe = form.RememberMe
	err = u.signIn(w, user)
	if err != nil {
		vd.SetAlert(err)
		u.LogInView.Render(w, r, vd)
		return
	}
	http.Redirect(w, r, "/user/articles", http.StatusFound)
}

// LogOut POST /logout
func (u *Users) LogOut(w http.ResponseWriter, r *http.Request) {
	// expire cookie
	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	user := context.User(r.Context())
	// error ignored - error is unlikely, and cannot recover here even if
	// there is one....
	token, _ := rand.RememberToken()
	user.Remember = token
	u.us.Update(user)
	http.Redirect(w, r, "/", http.StatusFound)
}

// InitiateReset POST /forgot
func (u *Users) InitiateReset(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form ResetPwForm
	vd.Yield = &form
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.ForgotPwView.Render(w, r, vd)
		return
	}

	token, err := u.us.InitiateReset(form.Email)
	if err != nil {
		vd.SetAlert(err)
		u.ForgotPwView.Render(w, r, vd)
		return
	}

	err = u.emailer.ResetPw(form.Email, token)
	if err != nil {
		vd.SetAlert(err)
		u.ForgotPwView.Render(w, r, vd)
		return
	}

	views.RedirectAlert(w, r, "/reset", http.StatusFound, views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Instructions for resetting your password have been emailed to you.",
	})

}

// ResetPW GET /reset
// Displays the result of the reset password form, has a method to prefill
// data.
func (u *Users) ResetPw(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form ResetPwForm
	vd.Yield = &form
	if err := parseURLParams(r, &form); err != nil {
		vd.SetAlert(err)
	}
	u.ResetPwView.Render(w, r, vd)
}

// CompleteReset POST /reset
func (u *Users) CompleteReset(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form ResetPwForm
	vd.Yield = &form
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.ResetPwView.Render(w, r, vd)
		return
	}

	user, err := u.us.CompleteReset(form.Token, form.Password)
	if err != nil {
		vd.SetAlert(err)
		u.ResetPwView.Render(w, r, vd)
		return
	}

	u.signIn(w, user)
	//TODO: REDIRECT
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
	articles, err := u.us.ArticlesByUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	vd.User = user
	vd.Yield = articles
	u.AllArticlesView.Render(w, r, vd)

}

// Profile returns user data and redenders a prefilled in and validated user
// form.
func (u *Users) Profile(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	user := context.User(r.Context())

	vd.User = user
	u.UserProfileView.Render(w, r, vd)
	// u.us.ByID(user.id)

}

// EditProfile ..
func (u *Users) EditProfile(w http.ResponseWriter, r *http.Request) {
	var form EditProfileForm
	user := context.User(r.Context())
	var vd = views.Data{}
	if err := parseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		// u.NewView.Render(w, r, vd)
		return
	}
	user.Bio = form.Bio

	err := u.us.Update(user)
	if err != nil {
		log.Fatal(err)
	}
	//TODO: redirect to profile page
	p := fmt.Sprintf("/user/%d", user.ID)
	http.Redirect(w, r, p, http.StatusFound)

}

// ShowUserProfile handles displaying public user information, it also passes
//
func (u *Users) ShowUserProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusNotFound)
		return
	}
	user, err := u.us.ByID(uint(id))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusNotFound)
		return
	}

	user.Articles, err = u.us.ArticlesByUser(user)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusNotFound)
		return
	}
	// TODO: think of a better way to do this:
	for i := range user.Articles {
		user.Articles[i].BodyHTML = template.HTML(user.Articles[i].Body)
		h := fmt.Sprintf("%v", user.Articles[i].BodyHTML)
		user.Articles[i].SnippedHTML = generateSnippet(h)
	}
	var vd = views.Data{}
	vd.Yield = user
	u.ProfileView.Render(w, r, vd)
}
