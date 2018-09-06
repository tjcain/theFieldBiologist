package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/tjcain/theFieldBiologist/context"
	"github.com/tjcain/theFieldBiologist/models"
	"github.com/tjcain/theFieldBiologist/views"
)

const (
	// ShowArticle is the named route for /article/:id
	ShowArticle = "show_article"
	// EditArticle is the named route for /article/:id/edit
	EditArticle = "edit_article"
	// lenSnippet dictates the length of snippet created
)

// ArticleForm holds data returned when creating or updating an article
type ArticleForm struct {
	Title string `schema:"title"`
	Body  string `schema:"body"`
}

// Articles represents the views associated with articles
type Articles struct {
	NewArticle      *views.View
	ShowFullArticle *views.View
	EditArticle     *views.View
	AllArticles     *views.View
	as              models.ArticleService
	r               *mux.Router
}

// NewArticles is used to create an Articles controller.
func NewArticles(as models.ArticleService, r *mux.Router) *Articles {
	return &Articles{
		NewArticle:      views.NewView("pages", "articles/new"),
		ShowFullArticle: views.NewView("pages", "articles/full"),
		EditArticle:     views.NewView("pages", "articles/edit"),
		AllArticles:     views.NewView("pages", "articles/all"),
		as:              as,
		r:               r,
	}
}

// Create - POST /articles
func (a *Articles) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form ArticleForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		a.NewArticle.Render(w, r, vd)
		return
	}
	user := context.User(r.Context())
	article := models.Article{
		UserID: user.ID,
		Title:  form.Title,
		Body:   []byte(form.Body),
	}
	if err := a.as.Create(&article); err != nil {
		vd.SetAlert(err)
		a.NewArticle.Render(w, r, vd)
		return
	}
	url, err := a.r.Get(ManageArticles).URL("id", strconv.Itoa(int(article.ID)))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)

}

// Update POST /galleries/:id/update
func (a *Articles) Update(w http.ResponseWriter, r *http.Request) {
	article, err := a.ArticleByID(w, r)
	if err != nil {
		// ArticleByID renders error
		return
	}
	user := context.User(r.Context())
	if article.UserID != user.ID {
		http.Error(w, "Permission denied", http.StatusForbidden)
		return
	}
	var vd views.Data
	vd.Yield = article
	var form ArticleForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		a.EditArticle.Render(w, r, vd)
		return
	}
	article.Title = form.Title
	article.Body = []byte(form.Body)
	err = a.as.Update(article)
	article.BodyHTML = template.HTML(form.Body)
	vd.Yield = article
	vd.Alert = &views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Article updated successfully",
	}
	a.EditArticle.Render(w, r, vd)
}

// Delete POST /articles/:id/delete
func (a *Articles) Delete(w http.ResponseWriter, r *http.Request) {
	article, err := a.ArticleByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if article.UserID != user.ID {
		http.Error(w, "You do not have permission to do that",
			http.StatusForbidden)
	}

	var vd views.Data
	err = a.as.Delete(article.ID)
	if err != nil {
		vd.SetAlert(err)
		vd.Yield = article
		a.EditArticle.Render(w, r, vd)
		return
	}
	url, err := a.r.Get(ManageArticles).URL("id", strconv.Itoa(int(article.ID)))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)

}

// Show GET /articles/:id
func (a *Articles) Show(w http.ResponseWriter, r *http.Request) {
	article, err := a.ArticleByID(w, r)
	if err != nil {
		// error is already rendered by ArticleByID
		return
	}
	var vd views.Data
	article.BodyHTML = template.HTML(article.Body)
	vd.Yield = article
	a.ShowFullArticle.Render(w, r, vd)
}

// ShowLatestArticles GET /articles
func (a *Articles) ShowLatestArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := a.as.GetAll()
	if err != nil {
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	// TODO: think of a better way to do this:
	for i, article := range articles {
		articles[i].BodyHTML = template.HTML(article.Body)
		h := fmt.Sprintf("%v", articles[i].BodyHTML)
		articles[i].SnippedHTML = generateSnippet(h)
	}

	var vd views.Data
	vd.Yield = articles
	a.AllArticles.Render(w, r, vd)
}

// Edit GET /gallereis/:id/edit
func (a *Articles) Edit(w http.ResponseWriter, r *http.Request) {
	article, err := a.ArticleByID(w, r)
	if err != nil {
		// ArticleByID renders error
		return
	}
	user := context.User(r.Context())
	if article.UserID != user.ID {
		http.Error(w, "Permission denied", http.StatusForbidden)
		return
	}
	article.BodyHTML = template.HTML(article.Body)
	var vd views.Data
	vd.Yield = article
	a.EditArticle.Render(w, r, vd)

}

// TODO: FIX THIS... SHOULD NOT BE MOVED TO MODELS
// ArticleByID retrieves an article by it's given id taken from the request
func (a *Articles) ArticleByID(w http.ResponseWriter,
	r *http.Request) (*models.Article, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusNotFound)
		return nil, err
	}
	article, err := a.as.ByID(uint(id))
	if err != nil {
		switch err {
		case models.ErrNotFound:
			http.Error(w, "Oops, article not found", http.StatusNotFound)
		default:
			http.Error(w, "Oops! Something went wrong.. my bad,",
				http.StatusInternalServerError)
		}
		return nil, err
	}
	return article, nil
}
