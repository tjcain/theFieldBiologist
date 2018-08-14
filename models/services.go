package models

import (
	"github.com/jinzhu/gorm"
	// db drivers
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Services ...
type Services struct {
	Article ArticleService
	User    UserService
	db      *gorm.DB
}

// NewServices opens a database conection, checks for errors, sets db log mode
// and constructs each service of the Service type.
func NewServices(connectionInfo string) (*Services, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &Services{
		User:    NewUserService(db),
		Article: NewArticleService(db),
		db:      db,
	}, nil
}

// DestructiveReset drops the user table and rebuilts it
func (s *Services) DestructiveReset() error {
	if err := s.db.DropTableIfExists(&User{}, Article{}).Error; err != nil {
		return err
	}
	return s.AutoMigrate()
}

// AutoMigrate will attempt to automatically migrate the users table
func (s *Services) AutoMigrate() error {
	if err := s.db.AutoMigrate(&User{}, Article{}).Error; err != nil {
		return err
	}
	return nil
}

// Close closes the UserService database connection
func (s *Services) Close() error {
	return s.db.Close()
}
