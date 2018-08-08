package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	homeTemplate *template.Template
)

func main() {

	r := mux.NewRouter()

	// Static Assets
	assetHandler := http.FileServer(http.Dir("./assets/"))
	assetHandler = http.StripPrefix("/assets/", assetHandler)
	r.PathPrefix("/assets/").Handler(assetHandler)

	// templates
	var err error
	homeTemplate, err = template.ParseFiles("views/home.gohtml")
	if err != nil {
		panic(err)
	}

	r.HandleFunc("/", home)

	fmt.Println("Starting server on port 3000")
	http.ListenAndServe(":3000", r)

}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := homeTemplate.Execute(w, nil); err != nil {
		panic(err)
	}
}
