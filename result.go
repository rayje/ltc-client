package main

import (
	"time"
)

type Result struct {
	Code      uint16
	Timestamp time.Time
	Latency   time.Duration
	BytesOut  uint64
	BytesIn   uint64
	Error     string
}

type Results []Result