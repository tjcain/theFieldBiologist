package controllers

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tjcain/theFieldBiologist/context"
	"github.com/tjcain/theFieldBiologist/models"
	"github.com/tjcain/theFieldBiologist/views"
)

// Admin represents a new admin view
type Admin struct {
	AdminDashView *views.View
	Articleview   *views.View
	us            models.UserService
	as            models.ArticleService
}

// AdminInfo ..
type AdminInfo struct {
	Users             uint
	DraftArticles     uint
	ReviewQueue       uint
	PublishedArticles uint
	Articles          []models.Article
}

// NewAdmin is used to create a new Admin controller.
// Any error in rendering templates will cause this function to panic.
func NewAdmin(us models.UserService, as models.ArticleService) *Admin {
	return &Admin{
		AdminDashView: views.NewView("pages", "admin/dashboard"),
		Articleview:   views.NewView("pages", "admin/articleview"),
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
	ArticlesForReview, err := a.as.ArticlesForReview()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	info.Users = userCount
	info.DraftArticles = DraftArticles
	info.ReviewQueue = ReviewQueue
	info.PublishedArticles = PublishedArticles
	info.Articles = ArticlesForReview
	var vd views.Data
	vd.Yield = info
	a.AdminDashView.Render(w, r, vd)
}

// ArticleView allows admin to edit, accept or reject a submitted article
func (a *Admin) ArticleView(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	article, err := a.as.ByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	article.BodyHTML = template.HTML(article.Body)
	var vd views.Data
	vd.Yield = article
	a.Articleview.Render(w, r, vd)
}

func (a *Admin) Accept(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	if !user.Admin {
		http.Error(w, "Permission denied", http.StatusForbidden)
		return
	}
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	article, err := a.as.ByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	article.Published = true
	article.Submitted = false
	err = a.as.Update(article)
	if err != nil {
		http.Error(w, "Ooops... something went wrong",
			http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/dashboard", http.StatusFound)
}

func (a *Admin) Reject(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	if !user.Admin {
		http.Error(w, "Permission denied", http.StatusForbidden)
		return
	}
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	article, err := a.as.ByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	article.Rejected = true
	article.Submitted = false
	err = a.as.Update(article)
	if err != nil {
		http.Error(w, "Ooops... something went wrong",
			http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	http.Redirect(w, r, "/admin/dashboard", http.StatusFound)
}
