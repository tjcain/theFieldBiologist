package controllers

import "github.com/tjcain/theFieldBiologist/views"

// Articles represents the views associated with articles
type Articles struct {
	NewArticle *views.View
}

// NewArticles is used to create an Articles controller.
func NewArticles() *Articles {
	return &Articles{
		NewArticle: views.NewView("pages", "articles/new"),
	}
}
