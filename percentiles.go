package main

import (
	"sort"
	"math"
	"time"
)

type Calculator struct {
	Results []Result
	Sum int64
}

func NewCalculator(results []Result, sum int64) Calculator {
	calc := Calculator{
		Results: results,
		Sum: sum,
	}

	return calc
}

func (c *Calculator) Sort() {
	sort.Sort(Results(c.Results))
}

func (c *Calculator) GetPercentile1(p float64) time.Duration {
	r := int64(float64(c.Sum) * p)
	
	for _, result := range c.Results {
		if result.Duration.Nanoseconds() > r {
			return result.Duration
		}
	}

	return time.Duration(0)
}

func (c *Calculator) GetPercentile2(p float64) time.Duration {
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