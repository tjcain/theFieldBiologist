package views

import (
	"log"

	"github.com/tjcain/theFieldBiologist/models"
)

const (
	// AlertLvlError indicates an error message
	AlertLvlError = "danger"
	// AlertLvlWarning indicates a warning message
	AlertLvlWarning = "warning"
	// AlertLvlInfo indiacates an informative message
	AlertLvlInfo = "info"
	// AlertLvlSuccess indicates positive reaction to a user action
	AlertLvlSuccess = "success"
	// AlertMsgGeneric is displayed to the user when an random error is
	// thrown by our backend.
	AlertMsgGeneric = "Oops, something went wrong. Please tray again, " +
		"and contact us if the problem persists."
)

// PublicError .. comment
type PublicError interface {
	error
	Public() string
}

// Data is the top level struct passed into views
type Data struct {
	Alert *Alert
	User  *models.User
	Yield interface{}
}

// SetAlert provides a method of creating an Alert with level Error and passing
// in an error.
func (d *Data) SetAlert(err error) {
	var msg string
	if pErr, ok := err.(PublicError); ok {
		msg = pErr.Public()
	} else {
		log.Println(err)
		msg = AlertMsgGeneric
	}
	d.Alert = &Alert{
		Level:   AlertLvlError,
		Message: msg,
	}
}

// AlertError is a helper method that makes it simple to set custom error alert
// messages
func (d *Data) AlertError(msg string) {
	d.Alert = &Alert{
		Level:   AlertLvlError,
		Message: msg,
	}
}

// Alert renders alert messages in templates
type Alert struct {
	Level   string
	Message string
}
