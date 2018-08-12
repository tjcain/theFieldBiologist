package models

import "github.com/jinzhu/gorm"

// Article represents data stored in the database for each article
type Article struct {
	gorm.Model
	UserID uint   `gorm:"not_null;index"`
	Title  string `gorm:"not_null"`
	Text   string
}

// ArticleService ...
type ArticleService interface {
	ArticleDB
}

// ArticleDB ...
type ArticleDB interface {
	Create(article *Article) error
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

type articleValidator struct {
	ArticleDB
}

var _ ArticleDB = &articleGorm{}

type articleGorm struct {
	db *gorm.DB
}

func (ag *articleGorm) Create(article *Article) error {
	return ag.db.Create(article).Error
}
