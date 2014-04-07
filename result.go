package main

import (
	"time"
	"sort"
)

type Result struct {
	Code      uint16
	Timestamp time.Time
	Duration  time.Duration
	BytesOut  uint64
	BytesIn   uint64
	Error     string
	ReadTime  time.Duration
	FanDuration []int64
	ReadTimes string
	Server string
}

type Results []Result

func (r Results) Sort() Results {
	sort.Sort(r)
	return r
}

func (r Results) Len() int           { return len(r) }
func (r Results) Less(i, j int) bool { return r[i].Duration < r[j].Duration}
func (r Results) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }