// +build !linux

package cbreaker

import "time"

// Ticker interface
type Ticker interface {
	Tick() int64
}

var startTick = time.Unix(0, 0)

type systemTicker struct{}

func (s *systemTicker) Tick() int64 {
	return time.Since(startTick).Nanoseconds()
}

// SystemTicker default ticker
var SystemTicker Ticker = &systemTicker{}
