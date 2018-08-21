package models

import (
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"
	"github.com/tjcain/theFieldBiologist/hash"
	"github.com/tjcain/theFieldBiologist/rand"
)

// UserDB is used to interact with the users database.
//
// single user queries:
// If the user is found, a nill error will be returned.
// If the user is not found ErrNotFound will be returned.
// If there is another error, we will return an error with
// more information about what went wrong. This may not be
// an error generated by the models package.
type UserDB interface {
	// Methods for querying for single users
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(rememberHash string) (*User, error)

	// Methods for altering users (CRUD operations)
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error

	// Methods for getting users articles
	ArticlesByUser(user *User) ([]Article, error)
}

// UserService is a set of methods used to manipulate and
// work with the user model
type UserService interface {
	// Authenticate will verify the provided email address and
	// password are correct. If they are correct, the user
	// corresponding to that email will be returned. Otherwise
	// You will receive either:
	// ErrNotFound, ErrPasswordInvalid, or another error if
	// something goes wrong.
	Authenticate(email, password string) (*User, error)
	UserDB
}

// User is a model of user details
type User struct {
	gorm.Model
	Name         string    `gorm:"not null"`
	Email        string    `gorm:"not null;unique_index"`
	Password     string    `gorm:"-"`
	PasswordHash string    `gorm:"not null"`
	Remember     string    `gorm:"-"`
	RememberHash string    `gorm:"not null;unique_index"`
	RememberMe   bool      `gorm:"-"`
	Articles     []Article `gorm:"foreignkey:UserID"`
}

// TODO: Delete before deployment
var _ UserService = &userService{}

// UserService provides an abstraction layer, and provides methods for
// querying, creating and updating users.
type userService struct {
	UserDB
	pepper string
}

// NewUserService instantiates a UserService with a connection to postgres db
func NewUserService(db *gorm.DB, pepper, hmacKey string) UserService {
	ug := &userGorm{db}
	hmac := hash.NewHMAC(hmacKey)
	uv := newUserValidator(ug, hmac, pepper)
	return &userService{
		UserDB: uv,
		pepper: pepper,
	}
}

// Authenticate is used to authenticate a user with the provided email and
// password.
// If the email address provided is invalid, this will return nil, ErrNotFound
// If the password provided is invalid, this will return nil, ErrPasswordInvalid
// If the email and password are both valid, thsi will return user, nil
// Otherwise, if any other error is encountered this will return nil, error
func (us *userService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(foundUser.PasswordHash), []byte(password+us.pepper))
	switch err {
	case nil:
		return foundUser, nil
	case bcrypt.ErrMismatchedHashAndPassword:
		return nil, ErrPasswordInvalid
	default:
		return nil, err
	}
}

// TODO: Delete me before deployment
var _ UserDB = &userGorm{}

// userGorm represents the database interaction layer and impliments the UserDB
// interface fully.
type userGorm struct {
	db *gorm.DB
}

// Create will create the provided user. Create expects a pre-validated user
func (ug *userGorm) Create(user *User) error {
	return ug.db.Create(user).Error
}

// Update will update a user. Update expectes a pre-validated user
func (ug *userGorm) Update(user *User) error {
	return ug.db.Save(user).Error
}

// Delete will delete a user with the provided id. Delete expects a pre
// sanitized id.
func (ug *userGorm) Delete(id uint) error {
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&user).Error
}

// ByID will look up a user with the provided ID.
// If the user is found, it will return a nil error.
// If no user is found, an ErrNotFound error will be returned.
// Any other error will result in an error being returned with more information
// about what went wrong. This may not be an error generated by the models
// package.
//
// As a general rule, any error other than ErrNotFound should probably result
// in a 500 error.
func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil

}

// ByEmail will look up a user with the provided email.
// If the user is found, it will return a nil error.
// If no user is found, an ErrNotFound error will be returned.
// Any other error will result in an error being returned with more information
// about what went wrong. This may not be an error generated by the models
// package.
//
// As a general rule, any error other than ErrNotFound should probably result
// in a 500 error.
func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email = ?", email)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, err
}

// ByRemember will look up a user with the provided remember token and return
// that user. This method expects a prehashed token.
func (ug *userGorm) ByRemember(rememberHash string) (*User, error) {
	var user User
	err := first(ug.db.Where("remember_hash = ?", rememberHash), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ArticlesByUser returns all articles linked to the provided user
func (ug *userGorm) ArticlesByUser(user *User) ([]Article, error) {
	var articles []Article
	err := ug.db.Model(&user).Related(&articles).Error
	if err != nil {
		return nil, err
	}
	return articles, nil
}

type userValFunc func(*User) error

func runUserValFns(user *User, fns ...userValFunc) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

// userValidator is our validation layer that validates
// and normalizes data before passing it on to the next
// UserDB in our interface chain.
type userValidator struct {
	UserDB
	hmac       hash.HMAC
	pepper     string
	emailRegex *regexp.Regexp
}

func newUserValidator(udb UserDB,
	hmac hash.HMAC, pepper string) *userValidator {
	return &userValidator{
		UserDB:     udb,
		hmac:       hmac,
		pepper:     pepper,
		emailRegex: regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`),
	}
}

// VALIDATION FUNCS

// bcryptPassword hashes user.Password using the bcypt package. It also will
// concatinate an app-wide pepper.  Salting is provided by bcrypt automatically
func (uv *userValidator) bcryptPassword(user *User) error {
	if user.Password == "" {
		return nil
	}
	pwBytes := []byte(user.Password + uv.pepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes,
		bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return nil
}

func (uv *userValidator) hmacRemember(user *User) error {
	if user.Remember == "" {
		return nil
	}
	user.RememberHash = uv.hmac.Hash(user.Remember)
	return nil
}

func (uv *userValidator) setRememberIfUnset(user *User) error {
	if user.Remember != "" {
		return nil
	}
	token, err := rand.RememberToken()
	if err != nil {
		return err
	}
	user.Remember = token
	return nil
}

func (uv *userValidator) idGreaterThan(n uint) userValFunc {
	return userValFunc(func(user *User) error {
		if user.ID <= n {
			return ErrIDInvalid
		}
		return nil
	})
}

func (uv *userValidator) normalizeEmail(user *User) error {
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)
	return nil
}

func (uv *userValidator) requireEmail(user *User) error {
	if user.Email == "" {
		return ErrEmailRequired
	}
	return nil
}

func (uv *userValidator) requireName(user *User) error {
	if user.Name == "" {
		return ErrNameRequired
	}
	return nil
}

func (uv *userValidator) emailFormat(user *User) error {
	if user.Email == "" {
		return nil
	}
	if !uv.emailRegex.MatchString(user.Email) {
		return ErrEmailInvalid
	}
	return nil
}

func (uv *userValidator) emailExistsCheck(user *User) error {
	existing, err := uv.ByEmail(user.Email)
	if err == ErrNotFound {
		return nil
	}
	if err != nil {
		return err
	}
	if user.ID != existing.ID {
		return ErrEmailTaken
	}
	return nil
}

func (uv *userValidator) passwordMinLength(user *User) error {
	if user.Password == "" {
		return nil
	}
	if len(user.Password) < 8 {
		return ErrPasswordTooShort
	}
	return nil
}

func (uv *userValidator) passwordRequired(user *User) error {
	if user.Password == "" {
		return ErrPasswordRequired
	}
	return nil
}

func (uv *userValidator) passwordHashRequired(user *User) error {
	if user.PasswordHash == "" {
		return ErrPasswordRequired
	}
	return nil
}

func (uv *userValidator) rememberMinBytes(user *User) error {
	if user.Remember == "" {
		return nil
	}
	n, err := rand.NBytes(user.Remember)
	if err != nil {
		return err
	}
	if n < 32 {
		return ErrRememberTooShort
	}
	return nil
}

func (uv *userValidator) rememberHashRequired(user *User) error {
	if user.RememberHash == "" {
		return ErrRememberRequired
	}
	return nil
}

// CRUD FUNCS

// Create will create the provided user and backfill data
// like the ID, CreatedAt, and UpdatedAt fields.
func (uv *userValidator) Create(user *User) error {
	err := runUserValFns(user,
		uv.requireName,
		uv.passwordRequired,
		uv.passwordMinLength,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.setRememberIfUnset,
		uv.rememberMinBytes,
		uv.hmacRemember,
		uv.rememberHashRequired,
		uv.requireEmail,
		uv.normalizeEmail,
		uv.emailFormat,
		uv.emailExistsCheck)
	if err != nil {
		return err
	}

	user.RememberHash = uv.hmac.Hash(user.Remember)
	return uv.UserDB.Create(user)
}

func (uv *userValidator) Update(user *User) error {
	err := runUserValFns(user,
		uv.requireName,
		uv.passwordMinLength,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.rememberMinBytes,
		uv.hmacRemember,
		uv.rememberHashRequired,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
		uv.emailExistsCheck)
	if err != nil {
		return err
	}
	return uv.UserDB.Update(user)
}

// Delete will delete the user with the provided ID, validating that there
// is no id of 0, which would delete the entire database...
func (uv *userValidator) Delete(id uint) error {
	var user User
	err := runUserValFns(&user, uv.idGreaterThan(0))
	if err != nil {
		return err
	}
	return uv.UserDB.Delete(id)
}

// ByRemember will hash the remember token and then call ug.ByRemember on the
// subusiquent UserDB layer
func (uv *userValidator) ByRemember(token string) (*User, error) {
	user := User{
		Remember: token,
	}
	if err := runUserValFns(&user, uv.hmacRemember); err != nil {
		return nil, err
	}
	return uv.UserDB.ByRemember(user.RememberHash)
}

// ByEmail will normalize an email address before passing
// it on to the database layer to perform the query.
func (uv *userValidator) ByEmail(email string) (*User, error) {
	user := User{
		Email: email,
	}
	err := runUserValFns(&user, uv.normalizeEmail)
	if err != nil {
		return nil, err
	}
	return uv.UserDB.ByEmail(user.Email)
}

// first will query using the provided gorm.DB. It will return the first item
// found and place it into dst.
//
// If nothing is found it will return ErrNotFound
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}
