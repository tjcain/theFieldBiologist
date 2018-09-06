package controllers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tjcain/theFieldBiologist/models"
	"github.com/tjcain/theFieldBiologist/views"
)

const (
	// PerPage dictates the amount of articles renderd to the featured article
	// section
	PerPage = 4
)

// HomePage holds the data that will be passed into the tempalte
type HomePage struct {
	EditorsPicks []models.Article
	// TODO: Impliment FeaturedArticles
	// FeaturedArticles []models.Article
	// TODO: Replace with  FeaturedArticle
	LatestArticles []models.Article
}

type Index struct {
	IndexPage *views.View
	as        models.ArticleService
	r         *mux.Router
}

func NewIndex(as models.ArticleService, r *mux.Router) *Index {
	return &Index{
		IndexPage: views.NewView("homePage", "static/home"),
		as:        as,
		r:         r,
	}
}

func (i *Index) Index(w http.ResponseWriter, r *http.Request) {
	articles, err := i.as.LatestArticles(PerPage)
	if err != nil {
		http.Error(w, "uhoh, something is wrong",
			http.StatusInternalServerError)
	}
	// TODO: think of a better way to do this:
	for i, article := range articles {
		articles[i].BodyHTML = template.HTML(article.Body)
		h := fmt.Sprintf("%v", articles[i].BodyHTML)
		articles[i].SnippedHTML = generateSnippet(h)
	}

	h := HomePage{
		LatestArticles: articles,
	}
	var vd views.Data
	vd.Yield = h
	i.IndexPage.Render(w, r, vd)
}
