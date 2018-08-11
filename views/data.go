package views

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

// Data is the top level struct passed into views
type Data struct {
	Alert *Alert
	Yield interface{}
}

// Alert renders alert messages in templates
type Alert struct {
	Level   string
	Message string
}
