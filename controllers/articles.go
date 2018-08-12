package controllers

import (
	"github.com/tjcain/theFieldBiologist/models"
	"github.com/tjcain/theFieldBiologist/views"
)

// Articles represents the views associated with articles
type Articles struct {
	NewArticle *views.View
	as         models.ArticleService
}

// NewArticles is used to create an Articles controller.
func NewArticles(as models.ArticleService) *Articles {
	return &Articles{
		NewArticle: views.NewView("pages", "articles/newarticle"),
		as:         as,
	}
}
