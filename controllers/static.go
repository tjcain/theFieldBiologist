package controllers

import "github.com/tjcain/theFieldBiologist/views"

// Static represents static page views
type Static struct {
	HomeView  *views.View
	AboutView *views.View
	// TODO: Rest of semi static pages:
	// Contact
	// News
}

func NewStatic() *Static {
	return &Static{
		HomeView:  views.NewView("homePage", "static/home"),
		AboutView: views.NewView("pages", "static/about"),
	}
}
