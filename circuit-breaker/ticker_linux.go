// +build linux

package cbreaker

import "time"

// Ticker interface
type Ticker interface {
	Tick() int64
}

type systemTicker struct{}

func (s *systemTicker) Tick() int64 {
	return time.Now().UnixNano()
}

// SystemTicker default ticker
var SystemTicker Ticker = &systemTicker{}
