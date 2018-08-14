package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tjcain/theFieldBiologist/controllers"
	"github.com/tjcain/theFieldBiologist/middleware"
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
	services, err := models.NewServices(psqlInfo)
	if err != nil {
		panic(err)
	}

	defer services.Close()
	services.AutoMigrate()

	r := mux.NewRouter()
	// Static Assets
	assetHandler := http.FileServer(http.Dir("./assets/"))
	assetHandler = http.StripPrefix("/assets/", assetHandler)
	r.PathPrefix("/assets/").Handler(assetHandler)

	// controllers
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User)
	articlesC := controllers.NewArticles(services.Article, r)

	// middleware
	requireUserMw := middleware.RequireUser{
		UserService: services.User,
	}

	// newArticle := requireUserMw.Apply(articlesC.NewArticle)
	// createArticle := requireUserMw.ApplyFn(articlesC.Create)

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

	// articles
	r.Handle("/article/new",
		requireUserMw.Apply(articlesC.NewArticle)).Methods("GET")
	r.HandleFunc("/article/new",
		requireUserMw.ApplyFn(articlesC.Create)).Methods("POST")
	// TODO: Move /articles to the user controller, this shows all of a single
	// users articles
	r.Handle("/articles",
		requireUserMw.ApplyFn(articlesC.ShowAllArticles)).Methods("GET")
	r.HandleFunc("/article/{id:[0-9]+}",
		articlesC.Show).Methods("GET").Name(controllers.ShowArticle)
	r.HandleFunc("/article/{id:[0-9]+}/edit",
		requireUserMw.ApplyFn(articlesC.Edit)).Methods("GET")
	r.HandleFunc("/article/{id:[0-9]+}/update",
		requireUserMw.ApplyFn(articlesC.Update)).Methods("POST")
	r.HandleFunc("/article/{id:[0-9]+}/delete",
		requireUserMw.ApplyFn(articlesC.Delete)).Methods("POST")

	// ifconfig | grep netmask
	fmt.Println("Listening on localhost:8080")
	// fmt.Println("Listening on local network:", devhelpers.LocalIP()+":8080")
	http.ListenAndServe(":8080", r)
}

// helper function that panics on error
func must(err error) {
	if err != nil {
		panic(err)
	}
}
