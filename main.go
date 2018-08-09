package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tjcain/theFieldBiologist/views"
)

var (
	homeView   *views.View
	aboutView  *views.View
	logInView  *views.View
	signUpView *views.View
)

func main() {

	r := mux.NewRouter()

	// Static Assets
	assetHandler := http.FileServer(http.Dir("./assets/"))
	assetHandler = http.StripPrefix("/assets/", assetHandler)
	r.PathPrefix("/assets/").Handler(assetHandler)

	// templates
	var err error
	homeView = views.NewView("homePage", "views/static/home.gohtml")
	if err != nil {
		panic(err)
	}

	aboutView = views.NewView("pages", "views/static/about.gohtml")
	if err != nil {
		panic(err)
	}

	logInView = views.NewView("pages", "views/users/login.gohtml")
	if err != nil {
		panic(err)
	}

	signUpView = views.NewView("pages", "views/users/signup.gohtml")
	if err != nil {
		panic(err)
	}

	r.HandleFunc("/", home)
	r.HandleFunc("/about", about)
	r.HandleFunc("/login", logIn)
	r.HandleFunc("/signup", signUp)

	fmt.Println("Starting server on port 3000")
	http.ListenAndServe(":3000", r)

}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(homeView.Render(w, nil))
}

func about(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(aboutView.Render(w, nil))
}

func logIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(logInView.Render(w, nil))
}

func signUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(signUpView.Render(w, nil))
}

// helper function that panics on error
func must(err error) {
	if err != nil {
		panic(err)
	}
}
