package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tjcain/theFieldBiologist/controllers"
	"github.com/tjcain/theFieldBiologist/models"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "postgres"
	dbName = "fieldbiologist_dev"
)

func main() {

	// Create a db connection
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname =%s "+
		"sslmode=disable", host, port, user, dbName)
	// Create a UserService
	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()
	us.AutoMigrate()

	r := mux.NewRouter()
	// Static Assets
	assetHandler := http.FileServer(http.Dir("./assets/"))
	assetHandler = http.StripPrefix("/assets/", assetHandler)
	r.PathPrefix("/assets/").Handler(assetHandler)

	// controllers
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(us)

	// handlers
	r.Handle("/", staticC.HomeView).Methods("GET")
	r.Handle("/about", staticC.AboutView).Methods("GET")

	// users
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/login", usersC.LogInView).Methods("GET")
	r.HandleFunc("/login", usersC.LogIn).Methods("POST")
	// cookietest is for dev only..
	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")

	fmt.Println("Starting server on port 3000")
	http.ListenAndServe(":3000", r)

}

// helper function that panics on error
func must(err error) {
	if err != nil {
		panic(err)
	}
}
