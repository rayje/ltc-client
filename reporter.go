package main

import (
	"fmt"
	"strings"
	"time"
)

func report(results []Result, name string) {
	calc := NewMetrics(results)
	calc.Sort()

	_25  := calc.GetPercentile(0.25)
	_75  := calc.GetPercentile(0.75)
	_95  := calc.GetPercentile(0.95)
	_99  := calc.GetPercentile(0.99)
	_999 := calc.GetPercentile(0.999)

	fmt.Println(strings.Repeat("=", 30))
	fmt.Println("Results -", name)
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("Requests:", calc.Total)

	fmt.Println("Latencies:")
	fmt.Printf("\tTotal:\t%s\t\n\n", calc.TotalRtt)

	fmt.Printf("\t0.25:\t%s\n", _25)
	fmt.Printf("\t0.75:\t%s\n", _75)
	fmt.Printf("\t0.95:\t%s\n", _95)
	fmt.Printf("\t0.99:\t%s\n", _99)
	fmt.Printf("\t0.999:\t%s\n\n", _999)

	fmt.Printf("\tmean:\t%s\n", time.Duration(calc.Mean))
	fmt.Printf("\tstd:\t%s\n", time.Duration(calc.StdDev()))
	fmt.Printf("\tmin:\t%s\n", calc.Min)
	fmt.Printf("\tmax:\t%s\n", calc.Max)
	fmt.Printf("\tiqr:\t%s\n", time.Duration(_75 - _25) )

	fmt.Println(strings.Repeat("=", 30))
}