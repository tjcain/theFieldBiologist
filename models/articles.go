package models

import (
	"html/template"

	"github.com/jinzhu/gorm"
)

// Article represents data stored in the database for each article
type Article struct {
	gorm.Model
	UserID   uint   `gorm:"not_null;index"`
	UserName string `gorm:"-"`
	Title    string `gorm:"not_null"`
	// stored as bytea in postrges... keeps things slim
	Body     []byte
	BodyHTML template.HTML `gorm:"-"`
}

// ArticleService ...
type ArticleService interface {
	ArticleDB
}

// ArticleDB ...
type ArticleDB interface {
	ByID(id uint) (*Article, error)
	Create(article *Article) error
	Update(article *Article) error
	Delete(id uint) error
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
	db := ag.db.Where("id = ?", id)
	err := first(db, &article)
	if err != nil {
		return nil, err
	}
	return &article, nil
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
