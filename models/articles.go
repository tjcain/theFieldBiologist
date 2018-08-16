package models

import (
	"html/template"

	"github.com/jinzhu/gorm"
)

// Article represents data stored in the database for each article
type Article struct {
	gorm.Model
	UserID uint   `gorm:"not_null;index"`
	Title  string `gorm:"not_null"`
	// Name represents the name from the users table, it is returned from a
	// join during database query
	Name string `gorm:"-"`
	// Body contains the raw HTML submitted from WYSIWYG editor, it is stored
	// as bytea in the databse. Note: this contains escaped html, for rendering
	// back to the webpage in a '.contents' div it needs to be cast to
	// template.HTML
	Body []byte
	// BodyHTML stores template.HTML type, which removes HTML escape characters,
	// it is needed to render the output from the WYSIWYG editor.
	BodyHTML template.HTML `gorm:"-"`
	// Snippet stores an N byte snippet of the text contained within the first
	// occurance of <p> </p> tags with ... appended.

	SnippedHTML template.HTML `gorm:"-"`
}

// ArticleService ...
type ArticleService interface {
	ArticleDB
}

// ArticleDB ...
type ArticleDB interface {
	ByID(id uint) (*Article, error)
	ByUserID(userID uint) ([]Article, error)
	GetAll() ([]Article, error)
	Create(article *Article) error
	Update(article *Article) error
	Delete(id uint) error
	LatestArticles(limit int) ([]Article, error)
}

type articleService struct {
	ArticleDB
}

// NewArticleService ...
func NewArticleService(db *gorm.DB) ArticleService {
	return &articleService{
		ArticleDB: &articleValidator{
			ArticleDB: &articleGorm{
				db: db,
			},
		},
	}
}

type articleValFn func(*Article) error

type articleValidator struct {
	ArticleDB
}

func runArticleValFns(article *Article, fns ...articleValFn) error {
	for _, fn := range fns {
		if err := fn(article); err != nil {
			return err
		}
	}
	return nil
}

// Create first validates article creation and then calls the create
// function from our database layer
func (av *articleValidator) Create(article *Article) error {
	err := runArticleValFns(article,
		av.userIDRequred,
		av.titleRequired)
	if err != nil {
		return err
	}
	return av.ArticleDB.Create(article)
}

// Update ...
func (av *articleValidator) Update(article *Article) error {
	err := runArticleValFns(article,
		av.userIDRequred,
		av.titleRequired)
	if err != nil {
		return err
	}
	return av.ArticleDB.Update(article)
}

// Delete ...
func (av *articleValidator) Delete(id uint) error {
	var article Article
	article.ID = id
	if err := runArticleValFns(&article, av.nonZeroID); err != nil {
		return err
	}
	return av.ArticleDB.Delete(article.ID)
}

func (av *articleValidator) userIDRequred(a *Article) error {
	if a.UserID <= 0 {
		return ErrUserIDRequired
	}
	return nil
}

func (av *articleValidator) titleRequired(a *Article) error {
	if a.Title == "" {
		return ErrTitleRequired
	}
	return nil
}

func (av *articleValidator) nonZeroID(a *Article) error {
	if a.ID <= 0 {
		return ErrIDInvalid
	}
	return nil
}

var _ ArticleDB = &articleGorm{}

type articleGorm struct {
	db *gorm.DB
}

func (ag *articleGorm) ByID(id uint) (*Article, error) {
	var article Article
	db := ag.db.Table("articles").Select("articles.*, users.name").
		Joins("join users on articles.user_id = users.id").
		Where("articles.id = ?", id)
	err := first(db, &article)
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (ag *articleGorm) ByUserID(userID uint) ([]Article, error) {
	var articles []Article
	db := ag.db.Table("articles").Select("articles.*, users.name").
		Joins("join users on articles.user_id = users.id").
		Where("user_id = ?", userID)
	if err := db.Find(&articles).Error; err != nil {
		return nil, err
	}
	return articles, nil
}

func (ag *articleGorm) GetAll() ([]Article, error) {
	var articles []Article
	db := ag.db.Table("articles").Select("articles.*, users.name").
		Joins("join users on articles.user_id = users.id")
	if err := db.Find(&articles).Error; err != nil {
		return nil, err
	}
	return articles, nil
}

func (ag *articleGorm) LatestArticles(limit int) ([]Article, error) {
	var articles []Article
	db := ag.db.Table("articles").Select("articles.*, users.name").
		Joins("join users on articles.user_id = users.id")
	if err := db.Limit(limit).Find(&articles).Error; err != nil {
		return nil, err
	}
	return articles, nil
}

// Create ...
func (ag *articleGorm) Create(article *Article) error {
	return ag.db.Create(article).Error
}

// Update ...
func (ag *articleGorm) Update(article *Article) error {
	return ag.db.Save(article).Error
}

// Delete ...
func (ag *articleGorm) Delete(id uint) error {
	article := Article{Model: gorm.Model{ID: id}}
	return ag.db.Delete(&article).Error
}
