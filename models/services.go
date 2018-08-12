package models

import "github.com/jinzhu/gorm"

// Services ...
type Services struct {
	Article ArticleService
	User    UserService
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
		Article: &articleGorm{},
	}, nil
}
