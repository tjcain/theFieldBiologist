package views

import (
	"log"
	"net/http"
	"time"

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

// RedirectAlert accepts all the normal params for an
// http.Redirect and performs a redirect, but only after
// persisting the provided alert in a cookie so that it can
// be displayed when the new page is loaded.
func RedirectAlert(w http.ResponseWriter, r *http.Request, urlStr string, code int, alert Alert) {
	persistAlert(w, alert)
	http.Redirect(w, r, urlStr, code)
}

func persistAlert(w http.ResponseWriter, alert Alert) {
	// We don't want alerts showing up days later. If the
	// user doesnt load the redirect in 5 minutes we will
	// just expire it.
	expiresAt := time.Now().Add(5 * time.Minute)
	lvl := http.Cookie{
		Name:     "alert_level",
		Value:    alert.Level,
		Expires:  expiresAt,
		HttpOnly: true,
	}
	msg := http.Cookie{
		Name:     "alert_message",
		Value:    alert.Message,
		Expires:  expiresAt,
		HttpOnly: true,
	}
	http.SetCookie(w, &lvl)
	http.SetCookie(w, &msg)
}

func clearAlert(w http.ResponseWriter) {
	lvl := http.Cookie{
		Name:     "alert_level",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}
	msg := http.Cookie{
		Name:     "alert_message",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}
	http.SetCookie(w, &lvl)
	http.SetCookie(w, &msg)
}

func getAlert(r *http.Request) *Alert {
	// If either cookie is missing we will assume the alert
	// is invalid and return nil
	lvl, err := r.Cookie("alert_level")
	if err != nil {
		return nil
	}
	msg, err := r.Cookie("alert_message")
	if err != nil {
		return nil
	}
	alert := Alert{
		Level:   lvl.Value,
		Message: msg.Value,
	}
	return &alert
}
