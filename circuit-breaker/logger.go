package cbreaker

// Logger interface
type Logger interface {
	Info(i string)
	Warn(title string, v interface{})
	Error(title string, v interface{})
}

var logger Logger

// SetDefaultLogger set default logger
func SetDefaultLogger(lg Logger) {
	logger = lg
}
