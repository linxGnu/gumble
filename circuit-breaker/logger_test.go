package cbreaker

import (
	"fmt"
	"testing"
)

var loggedInfo string

var loggedWarn string

var loggedError string

type fakeLogger struct{}

func (f *fakeLogger) Info(i string) {
	loggedInfo = i
}

func (f *fakeLogger) Warn(title string, v interface{}) {
	loggedWarn = title + fmt.Sprintf("_%v", v)
}

func (f *fakeLogger) Error(title string, v interface{}) {
	loggedError = title + fmt.Sprintf("_%v", v)
}

func TestSetLogger(t *testing.T) {
	SetDefaultLogger(&fakeLogger{})

	if logger.Info("info"); loggedInfo != "info" {
		t.FailNow()
	}

	if logger.Warn("warn", "test"); loggedWarn != "warn_test" {
		t.FailNow()
	}

	if logger.Error("error", "test"); loggedError != "error_test" {
		t.FailNow()
	}

	if SetDefaultLogger(nil); logger != nil {
		t.FailNow()
	}
}
