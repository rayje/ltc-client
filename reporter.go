package main

import (
	"fmt"
	"strings"
	"time"
	"os"
)

func report(results []Result, name string, displayFan bool) {
	calc := NewMetrics(results)
	calc.Sort()

	_25m  := calc.GetPercentile(0.25)
	_75m  := calc.GetPercentile(0.75)
	_95m  := calc.GetPercentile(0.95)
	_99m  := calc.GetPercentile(0.99)
	_999m := calc.GetPercentile(0.999)

	var fan FanMetrics
	var _25f time.Duration
	var _75f time.Duration
	var _95f time.Duration
	var _99f time.Duration
	var _999f time.Duration
	var _minf time.Duration
	var _maxf time.Duration
	var _meanf float64
	var _stddevf float64
	if displayFan {
		fan = NewFanMetrics(results)
		fan.Sort()
		_25f  = fan.GetPercentile(0.25)
		_75f  = fan.GetPercentile(0.75)
		_95f  = fan.GetPercentile(0.95)
		_99f  = fan.GetPercentile(0.99)
		_999f = fan.GetPercentile(0.999)
		_minf = fan.Min
		_maxf = fan.Max
		_stddevf = fan.StdDev()
		_meanf = fan.Mean
	}

	fmt.Println(strings.Repeat("=", 30))
	fmt.Println("Results -", name)
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("Requests:", calc.Total)
	if displayFan {
		fmt.Println("FanRequests:", fan.Total)
	}

	fmt.Println("Latencies:")
	fmt.Printf("\tTotal:     %s\n", calc.TotalRtt)
	if displayFan {
		fmt.Printf("\tTotalFan:  %s\n\n", fan.TotalRtt)
	}

	fmt.Printf("\t     \t   Client")
	fmt.Printf("\t    Fan\n")
	fmt.Printf("\t\t" + strings.Repeat("-", 12))
	fmt.Printf("\t" + strings.Repeat("-", 12) + "\n")

	fmt.Printf("\t0.25:\t%s\t%s\n", _25m, _25f)
	fmt.Printf("\t0.75:\t%s\t%s\n", _75m, _75f)
	fmt.Printf("\t0.95:\t%s\t%s\n", _95m, _95f)
	fmt.Printf("\t0.99:\t%s\t%s\n", _99m, _99f)
	fmt.Printf("\t0.999:\t%s\t%s\n\n", _999m, _999f)

	fmt.Printf("\tmean:\t%s\t%s\n", time.Duration(calc.Mean), time.Duration(_meanf))
	fmt.Printf("\tstd:\t%s\t%s\n", time.Duration(calc.StdDev()), time.Duration(_stddevf))
	fmt.Printf("\tmin:\t%s\t%s\n", calc.Min, _minf)
	fmt.Printf("\tmax:\t%s\t%s\n", calc.Max, _maxf)
	fmt.Printf("\tiqr:\t%s\t%s\n", time.Duration(_75m - _25m), time.Duration(_75f - _25f) )

	fmt.Println(strings.Repeat("=", 30))
}

func dumpToFile(results []Result) {
    f, err := os.Create("results.txt")
    check(err)

    defer func() {
        if err := f.Close(); err != nil {
            panic(err)
        }
    }()

    for _, result := range results {
    	_, err := f.WriteString(fmt.Sprintf("%d", result.Duration.Nanoseconds()) + "\n")
    	check(err)
    }
}

func dumpFanToFile(results []Result) {
    f, err := os.Create("fan-results.txt")
    check(err)

    defer func() {
        if err := f.Close(); err != nil {
            panic(err)
        }
    }()

    for _, result := range results {
    	durations := make([]string, len(result.FanDuration))
		for i := 0; i < len(result.FanDuration); i++ {
			durations[i] = fmt.Sprintf("%d", result.FanDuration[i])
		}
    	_, err := f.WriteString(strings.Join(durations, ",") + "\n")
    	check(err)
    }
}

func dumpReadToFile(results []Result) {
    f, err := os.Create("results-read.txt")
    check(err)

    defer func() {
        if err := f.Close(); err != nil {
            panic(err)
        }
    }()

    var duration time.Duration

    for _, result := range results {
    	times := strings.Split(result.ReadTimes, ",")
    	durations := make([]string, len(times))

    	for i := 0; i < len(times); i++ {
	    	duration, err = time.ParseDuration(times[i])
	    	if err != nil {
	    		fmt.Println(err)
	    		fmt.Println("index:", i, times[i])
	    		duration = time.Duration(0)
	    	}
	    	durations[i] = fmt.Sprintf("%d", duration.Nanoseconds())
	    }

    	_, err := f.WriteString(strings.Join(durations,",") + "\n")
    	check(err)
    }
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}