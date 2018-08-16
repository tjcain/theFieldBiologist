package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/tjcain/theFieldBiologist/controllers"
	"github.com/tjcain/theFieldBiologist/devhelpers"
	"github.com/tjcain/theFieldBiologist/middleware"
	"github.com/tjcain/theFieldBiologist/models"
	"github.com/tjcain/theFieldBiologist/rand"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "postgres"
	dbName = "fieldbiologist_dev"
)

func main() {
	boolPtr := flag.Bool("prod", false, "Provide this flag in production. "+
		"This ensures a .config file exists before the application starts.")
	flag.Parse()
	// Load config file or return defaults
	cfg := LoadConfig(*boolPtr)
	dbCfg := cfg.Database
	// Create a Services
	services, err := models.NewServices(
		models.WithGorm(dbCfg.Dialect(), dbCfg.connectionInfo()),
		models.WithLogMode(!cfg.IsProd()),
		models.WithUser(cfg.Pepper, cfg.HMACKey),
		models.WithArticle(),
	)

	// services, err := models.NewServices(dbCfg.Dialect(), dbCfg.connectionInfo())
	// if err != nil {
	// 	panic(err)
	// }
	defer services.Close()
	services.AutoMigrate()

	// Cross-Site Request Forgery Protection
	b, err := rand.Bytes(32)
	if err != nil {
		panic(err)
	}
	csrfMw := csrf.Protect(b, csrf.Secure(cfg.IsProd()))

	r := mux.NewRouter()
	// Static Assets
	assetHandler := http.FileServer(http.Dir("./assets/"))
	assetHandler = http.StripPrefix("/assets/", assetHandler)
	r.PathPrefix("/assets/").Handler(assetHandler)

	// controllers
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User)
	articlesC := controllers.NewArticles(services.Article, r)
	indexC := controllers.NewIndex(services.Article, r)

	// middleware
	userMw := middleware.User{
		UserService: services.User,
	}
	requireUserMw := middleware.RequireUser{}

	// newArticle := requireUserMw.Apply(articlesC.NewArticle)
	// createArticle := requireUserMw.ApplyFn(articlesC.Create)

	// handlers
	r.HandleFunc("/", indexC.Index).Methods("GET")
	r.Handle("/about", staticC.AboutView).Methods("GET")

	// users
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/login", usersC.LogInView).Methods("GET")
	r.HandleFunc("/login", usersC.LogIn).Methods("POST")
	r.Handle("/logout", requireUserMw.ApplyFn(usersC.LogOut)).Methods("POST")
	r.Handle("/user/articles",
		requireUserMw.ApplyFn(usersC.ShowAllArticles)).Methods("GET").
		Name(controllers.ManageArticles)

	// cookietest is for dev only..
	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")

	// articles
	r.Handle("/article/new",
		requireUserMw.Apply(articlesC.NewArticle)).Methods("GET")
	r.HandleFunc("/article/new",
		requireUserMw.ApplyFn(articlesC.Create)).Methods("POST")
	r.HandleFunc("/articles", articlesC.ShowLatestArticles).Methods("GET")
	r.HandleFunc("/article/{id:[0-9]+}",
		articlesC.Show).Methods("GET").Name(controllers.ShowArticle)
	r.HandleFunc("/article/{id:[0-9]+}/edit",
		requireUserMw.ApplyFn(articlesC.Edit)).Methods("GET").
		Name(controllers.EditArticle)
	r.HandleFunc("/article/{id:[0-9]+}/update",
		requireUserMw.ApplyFn(articlesC.Update)).Methods("POST")
	r.HandleFunc("/article/{id:[0-9]+}/delete",
		requireUserMw.ApplyFn(articlesC.Delete)).Methods("POST")

	// startup
	fmt.Printf("Listening on localhost:%d\n", cfg.Port)
	// local network
	if !cfg.IsProd() {
		fmt.Printf("Listening on local network %s:%d\n", devhelpers.LocalIP(), cfg.Port)
	}
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), csrfMw(userMw.Apply(r)))
}

// helper function that panics on error
func must(err error) {
	if err != nil {
		panic(err)
	}
}
