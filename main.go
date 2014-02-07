package main

import (
	"fmt"
	"time"
)

func main() {
	config := getConfig()
	statsd := NewStatsd(&config)
	requestor, err := NewRequestor(&config, statsd)
	if err != nil {
		fmt.Println(err)
		return
	}

	var complete chan string
	reportsInterval := runReportsInterval(&config);
	if reportsInterval {
		complete = make(chan string)
		go intervalReports(&requestor, &config, complete)
	}

	results, err := requestor.MakeRequest()
	if err != nil {
		fmt.Println(err)
		return
	}

	if reportsInterval {
		complete <- "complete"
	}
	report(results, "final", config.Fan)
	dumpToFile(results)

	if config.Fan {
		dumpFanToFile(results)
		dumpReadToFile(results)
	}
}

func intervalReports(requestor *Requestor, config *Config, complete chan string) {
	var results Results
	runReport := time.Tick(config.ReportInterval)
	fmt.Println("ReportInterval:", config.ReportInterval)

	for {
		select {
		case <-complete:
		    return
		case <-runReport:
			results = requestor.GetResults()
			report(results, requestor.Url(), config.Fan)
		}
	}
}

func runReportsInterval(config *Config) bool {
	// Do no run reports on interval when set to 0
	if config.ReportInterval == 0 {
		return false
	}

	if config.ReportInterval >= config.Duration {
		return false
	}

	return true
}