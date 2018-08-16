package controllers

import (
	"github.com/tjcain/theFieldBiologist/views"
)

// Static represents static page views
type Static struct {
	AboutView *views.View
	// TODO: Rest of semi static pages:
	// Contact
	// News
}

// NewStatic returns a static controller
func NewStatic() *Static {
	return &Static{
		AboutView: views.NewView("pages", "static/about"),
	}
}
