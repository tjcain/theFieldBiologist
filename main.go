package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tjcain/theFieldBiologist/controllers"
	"github.com/tjcain/theFieldBiologist/views"
)

var (
	homeView  *views.View
	aboutView *views.View
)

func main() {

	r := mux.NewRouter()

	// Static Assets
	assetHandler := http.FileServer(http.Dir("./assets/"))
	assetHandler = http.StripPrefix("/assets/", assetHandler)
	r.PathPrefix("/assets/").Handler(assetHandler)

	homeView = views.NewView("homePage", "views/static/home.gohtml")
	aboutView = views.NewView("pages", "views/static/about.gohtml")

	// users controller
	usersC := controllers.NewUsers()

	r.HandleFunc("/", home).Methods("GET")
	r.HandleFunc("/about", about).Methods("GET")
	// r.HandleFunc("/login", logIn).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")

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

// helper function that panics on error
func must(err error) {
	if err != nil {
		panic(err)
	}
}
