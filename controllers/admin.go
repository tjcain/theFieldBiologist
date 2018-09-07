package controllers

import (
	"net/http"

	"github.com/tjcain/theFieldBiologist/models"
	"github.com/tjcain/theFieldBiologist/views"
)

// Admin represents a new admin view
type Admin struct {
	AdminDashView *views.View
	us            models.UserService
	as            models.ArticleService
}

// AdminInfo ..
type AdminInfo struct {
	Users             uint
	DraftArticles     uint
	ReviewQueue       uint
	PublishedArticles uint
}

// NewAdmin is used to create a new Admin controller.
// Any error in rendering templates will cause this function to panic.
func NewAdmin(us models.UserService, as models.ArticleService) *Admin {
	return &Admin{
		AdminDashView: views.NewView("pages", "admin/dashboard"),
		us:            us,
		as:            as,
	}
}

// Dashboard is the method that displays the admin dashboard
func (a *Admin) Dashboard(w http.ResponseWriter, r *http.Request) {
	var info AdminInfo
	userCount, err := a.us.UsersCount()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	DraftArticles, err := a.as.DraftArticles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ReviewQueue, err := a.as.ReviewQueue()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	PublishedArticles, err := a.as.PublishedArticles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	info.Users = userCount
	info.DraftArticles = DraftArticles
	info.ReviewQueue = ReviewQueue
	info.PublishedArticles = PublishedArticles
	var vd views.Data
	vd.Yield = info
	a.AdminDashView.Render(w, r, vd)

}
