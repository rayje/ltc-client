package main

import (
	"sort"
	"math"
	"time"
)

type Metrics struct {
	Results []Result
	Min time.Duration
	Max time.Duration
	TotalRtt time.Duration
	Mean float64
	Total float64
}

func NewMetrics(results []Result) Metrics {
	totalRtt := time.Duration(0)
	min := time.Hour
	max := time.Duration(0)

	for _, result := range results {
		totalRtt += result.Duration

		if result.Duration > max {
			max = result.Duration
		}

		if result.Duration < min {
			min = result.Duration
		}
	}

	total := float64(len(results))
	mean := float64(totalRtt)/total

	metric := Metrics{
		Results: results,
		Min: min,
		Max: max,
		TotalRtt: totalRtt,
		Mean: mean,
		Total: total,
	}

	return metric
}

func (c *Metrics) Sort() {
	sort.Sort(Results(c.Results))
}

func (c *Metrics) StdDev() float64 {
	var diffs float64
	m := float64(time.Duration(c.Mean).Nanoseconds())

	for _, result := range c.Results {
		diffs += math.Pow(float64(result.Duration.Nanoseconds()) - m, 2)
	}

	variance := diffs / c.Total
	stdDev := math.Sqrt(variance)

	return stdDev
}

func (c *Metrics) GetPercentile(p float64) time.Duration {
	var percentile float64

	r := (float64(len(c.Results)) + 1.0) * p
	ir, fr := math.Modf(r)

	v1 := float64(c.Results[int(ir)-1].Duration.Nanoseconds())

	if fr > 0.0 && ir < float64(len(c.Results)){
		v2 := float64(c.Results[int(ir)].Duration.Nanoseconds())
		percentile = (v2 - v1) * fr + v1
	} else {
		percentile = v1
	}

	return time.Duration(percentile)
}

