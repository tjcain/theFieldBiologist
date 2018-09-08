package controllers

import (
	"github.com/tjcain/theFieldBiologist/views"
)

// Static represents static page views
type Static struct {
	AboutView   *views.View
	ContactView *views.View
	PrivacyView *views.View
	// TODO: Rest of semi static pages:
	// News
}

// NewStatic returns a static controller
func NewStatic() *Static {
	return &Static{
		AboutView:   views.NewView("pages", "static/about"),
		ContactView: views.NewView("pages", "static/contact"),
		PrivacyView: views.NewView("pages", "static/privacy"),
	}
}
