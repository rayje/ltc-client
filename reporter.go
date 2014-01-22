package main

import (
	"fmt"
	"github.com/bmizerany/perks/quantile"
	"strings"
	"time"
	"math"
)

func report(results []Result, name string) {
	quants := quantile.NewTargeted(0.95, 0.99)

	total := float64(len(results))
	totalRtt := time.Duration(0)
	min := time.Hour
	max := time.Duration(0)

	for _, result := range results {
		quants.Insert(float64(result.Duration))
		totalRtt += result.Duration

		if result.Duration > max {
			max = result.Duration
		}

		if result.Duration < min {
			min = result.Duration
		}
	}

	mean := float64(totalRtt)/total
	stdDev := getDistributions(total, mean, results)

	_75 := quants.Query(0.75)
	_25 := quants.Query(0.25)

	fmt.Println(strings.Repeat("=", 30))
	fmt.Println("Results -", name)
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("Requests:", total)

	fmt.Println("Latencies:")
	fmt.Printf("\tTotal:\t%s\t\n\n", totalRtt)

	fmt.Printf("\t0.25:\t%s\n", time.Duration(_25))
	fmt.Printf("\t0.75:\t%s\n", time.Duration(_75))
	fmt.Printf("\t0.95:\t%s\n", time.Duration(quants.Query(0.95)))
	fmt.Printf("\t0.99:\t%s\n", time.Duration(quants.Query(0.99)))
	fmt.Printf("\t0.999:\t%s\n\n", time.Duration(quants.Query(0.999)))

	fmt.Printf("\tmean:\t%s\n", time.Duration(mean))
	fmt.Printf("\tstd:\t%s\n", time.Duration(stdDev))
	fmt.Printf("\tmin:\t%s\n", min)
	fmt.Printf("\tmax:\t%s\n", max)
	fmt.Printf("\tiqr:\t%s\n", time.Duration(_75 - _25) )
	fmt.Println(strings.Repeat("=", 30))
}

func getDistributions(numResults float64, mean float64, results []Result) float64 {
	var diffs float64
	m := float64(time.Duration(mean).Nanoseconds())

	for _, result := range results {
		diffs += math.Pow(float64(result.Duration.Nanoseconds()) - m, 2)
	}

	variance := diffs / numResults
	stdDev := math.Sqrt(variance)

	return stdDev
}