package models

import "strings"

const (
	// PUBLIC ERRORS - displayed to end user

	// ErrNotFound is returned when a resource cannot be found
	ErrNotFound modelError = "models: resource not found"
	// ErrPasswordInvalid is returned when a user attemptes to log in using an
	// incorrect password and cannot be authenticated
	ErrPasswordInvalid modelError = "models: incorrect password provided"
	// ErrNameRequired is returned when the name address field is empty when
	// attempting to create a user
	ErrNameRequired modelError = "models: a name is required"
	// ErrEmailRequired is returned when an email address field is empty when
	// attempting to create a user
	ErrEmailRequired modelError = "models: an email address is required"
	// ErrEmailInvalid is returned when a provided email address
	// does not pass validation
	ErrEmailInvalid modelError = "models: email address is invalid"
	// ErrEmailTaken is returned when an update or create call is attempted
	// on an email address that is already in the database.
	ErrEmailTaken modelError = "models: email address is already taken"
	// ErrPasswordTooShort is returned with a user attempts to set a password
	// that is less than 8 characters
	ErrPasswordTooShort modelError = "models: password must be at least 8" +
		" characters long"
	// ErrPasswordRequired is returned when a create is attempted with a null
	// password field.
	ErrPasswordRequired modelError = "models: password is required"
	// ErrTitleRequired is returned when an attempt to create an asset is made
	// without a title
	ErrTitleRequired modelError = "models: title is required"

	//PRIVATE ERRORS - not displayed to end user

	// ErrIDInvalid is returned when an invalid ID is provied to a method
	ErrIDInvalid privateError = "models: ID provided was invalid"
	// ErrRememberRequired is returned when a create or update is attempted
	// without a user remember token hash
	ErrRememberRequired privateError = "models: remember token is required"
	// ErrRememberTooShort is retunred when a remember token generated is less
	// than 32 bytes
	ErrRememberTooShort privateError = "models: remember token must be atleast" +
		" 32 bytes"
	// ErrUserIDRequired is returned when a user id is not provided during the
	// creation of an article or other user linked asset
	ErrUserIDRequired privateError = "models: user ID is required"
)

type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	split := strings.Split(s, " ")
	split[0] = strings.Title(split[0])
	return strings.Join(split, " ")
}

type privateError string

func (e privateError) Error() string {
	return string(e)
}
