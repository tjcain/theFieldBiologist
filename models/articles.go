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

type articleGorm struct {
	db *gorm.DB
}

func (gg *articleGorm) Create(article *Article) error {
	// TODO: Impliment this
	return nil
}
