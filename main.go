package main

import (
	"fmt"
)

func main() {
	config := getConfig()
	statsd := NewStatsd(&config)
	requestor := NewRequestor(&config)

	results, err := requestor.makeRequest(statsd)
	if err != nil {
		fmt.Println(err)
		return
	}

	report(results)
}
