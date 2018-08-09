package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tjcain/theFieldBiologist/controllers"
)

func main() {
	r := mux.NewRouter()
	// Static Assets
	assetHandler := http.FileServer(http.Dir("./assets/"))
	assetHandler = http.StripPrefix("/assets/", assetHandler)
	r.PathPrefix("/assets/").Handler(assetHandler)

	// controllers
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers()

	// handlers
	r.Handle("/", staticC.HomeView).Methods("GET")
	r.Handle("/about", staticC.AboutView).Methods("GET")
	// r.HandleFunc("/login", logIn).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")

	fmt.Println("Starting server on port 3000")
	http.ListenAndServe(":3000", r)

}

// helper function that panics on error
func must(err error) {
	if err != nil {
		panic(err)
	}
}
