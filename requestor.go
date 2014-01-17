package main

import (
	"io/ioutil"
	"net/http"
	"time"
)

var client = &http.Client{}

func NewRequest(url string) (http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Client", "Header")

	return *req, err
}

func makeRequest(url string, rate uint64, duration time.Duration) (Results, error) {
	total := rate * uint64(duration.Seconds())
	res := make(chan Result, total)
	results := make(Results, total)

	req, err := NewRequest(url)
	if err != nil {
		return nil, err
	}

	go runRequests(rate, &req, res, total)

	for i := 0; i < cap(res); i++ {
		results[i] = <- res
	}
	close(res)

	return results, nil
}

func runRequests(rate uint64, req *http.Request, res chan Result, total uint64) {
	throttle := time.Tick(time.Duration(1e9 / rate))

	for i := 0; uint64(i) < total; i++ {
		<-throttle
		go runRequest(req, res)
	}
}

func runRequest(req *http.Request, res chan Result) {
	start := time.Now()
	r, err := client.Do(req)

	result := Result {
		Timestamp: start,
		RTT:   time.Since(start),
		BytesOut:  uint64(req.ContentLength),
	}

	if err != nil {
		result.Error = err.Error()
	} else {
		result.Code = uint16(r.StatusCode)
		if body, err := ioutil.ReadAll(r.Body); err != nil {
			if result.Code < 200 || result.Code >= 300 {
				result.Error = string(body)
			}
		} else {
			result.BytesIn = uint64(len(body))
		}
	}

	res <- result
}