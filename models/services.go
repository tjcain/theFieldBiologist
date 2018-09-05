package models

import (
	"github.com/jinzhu/gorm"
	// postgres gorm drivers
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Services holds values required to init services
type Services struct {
	Article ArticleService
	User    UserService
	db      *gorm.DB
}

// NewServices accepts a list of config functions to run. Each function must
// accept a pointer to the current Services object.
func NewServices(cfgs ...ServicesConfig) (*Services, error) {
	var s Services
	for _, cfg := range cfgs {
		if err := cfg(&s); err != nil {
			return nil, err
		}
	}
	return &s, nil
}

// // NewServices opens a database conection, checks for errors, sets db log mode
// // and constructs each service of the Service type.
// func NewServices(dialect, connectionInfo string) (*Services, error) {
// 	db, err := gorm.Open(dialect, connectionInfo)
// 	if err != nil {
// 		return nil, err
// 	}
// 	db.LogMode(true)
// 	return &Services{
// 		User:    NewUserService(db),
// 		Article: NewArticleService(db),
// 		db:      db,
// 	}, nil
// }

// DestructiveReset drops the user table and rebuilts it
func (s *Services) DestructiveReset() error {
	if err := s.db.DropTableIfExists(&User{}, Article{}).Error; err != nil {
		return err
	}
	return s.AutoMigrate()
}

// AutoMigrate will attempt to automatically migrate the users and article table
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

type ServicesConfig func(*Services) error

func WithGorm(dialect, connectionInfo string) ServicesConfig {
	return func(s *Services) error {
		db, err := gorm.Open(dialect, connectionInfo)
		if err != nil {
			return err
		}
		s.db = db
		return nil
	}
}

func WithLogMode(mode bool) ServicesConfig {
	return func(s *Services) error {
		s.db.LogMode(mode)
		return nil
	}
}

func WithUser(pepper, hmacKey string) ServicesConfig {
	return func(s *Services) error {
		s.User = NewUserService(s.db, pepper, hmacKey)
		return nil
	}
}

func WithArticle() ServicesConfig {
	return func(s *Services) error {
		s.Article = NewArticleService(s.db)
		return nil
	}
}

// Future service configs will go here.
