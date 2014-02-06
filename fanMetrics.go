package main

import (
	"sort"
	"math"
	"time"
)

type FanMetrics struct {
	Results []int64
	Min time.Duration
	Max time.Duration
	TotalRtt time.Duration
	Mean float64
	Total float64
}

type int64arr []int64
func (a int64arr) Len() int { return len(a) }
func (a int64arr) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a int64arr) Less(i, j int) bool { return a[i] < a[j] }


func NewFanMetrics(results []Result) FanMetrics {
	totalRtt := time.Duration(0)
	min := time.Hour
	max := time.Duration(0)
	durations := []int64{}

	for _, result := range results {
		durations = append(durations, result.FanDuration...)

		for _, durationValue := range result.FanDuration {
			duration := time.Duration(durationValue)
			totalRtt += duration

			if duration > max {
				max = duration
			}

			if duration < min {
				min = duration
			}
		}
	}

	total := float64(len(durations))
	mean := float64(totalRtt)/total

	metric := FanMetrics{
		Results: durations,
		Min: min,
		Max: max,
		TotalRtt: totalRtt,
		Mean: mean,
		Total: total,
	}

	return metric
}

func (c *FanMetrics) Sort() {
	sort.Sort(int64arr(c.Results))
}

func (c *FanMetrics) StdDev() float64 {
	var diffs float64
	m := float64(time.Duration(c.Mean).Nanoseconds())

	for _, result := range c.Results {
		diffs += math.Pow(float64(result) - m, 2)
	}

	variance := diffs / c.Total
	stdDev := math.Sqrt(variance)

	return stdDev
}

func (c *FanMetrics) GetPercentile(p float64) time.Duration {
	var percentile float64

	r := (float64(len(c.Results)) + 1.0) * p
	ir, fr := math.Modf(r)

	v1 := float64(c.Results[int(ir)-1])

	if fr > 0.0 && ir < float64(len(c.Results)){
		v2 := float64(c.Results[int(ir)])
		percentile = (v2 - v1) * fr + v1
	} else {
		percentile = v1
	}

	return time.Duration(percentile)
}