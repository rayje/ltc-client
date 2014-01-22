package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Requestor struct {
	Rate     uint64
	Duration time.Duration
	EndPoint EndPoint
}

type Result struct {
	Code      uint16
	Timestamp time.Time
	Duration  time.Duration
	BytesOut  uint64
	BytesIn   uint64
	Error     string
}

type Results []Result

var client = &http.Client{}

func NewRequest(url string) (http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Client", "Header")

	return *req, err
}

func NewRequestor(config *Config) Requestor {
	requestor := &Requestor{
		Rate:     config.Rate,
		Duration: config.Duration,
		EndPoint: config.Endpoint,
	}

	return *requestor
}

func (r *Requestor) Url() string {
	return fmt.Sprintf("http://%s:%s/%s", r.EndPoint.Host, r.EndPoint.Port, r.EndPoint.Route)
}

func (r *Requestor) makeRequest(statsd StatsdClient) (Results, error) {
	total := r.Rate * uint64(r.Duration.Seconds())
	res := make(chan Result, total)
	results := make(Results, total)

	req, err := NewRequest(r.Url())
	if err != nil {
		return nil, err
	}

	go runRequests(r.Rate, &req, res, total)

	for i := 0; i < cap(res); i++ {
		results[i] = <-res
		statsd.Timing(results[i].Duration)
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

	result := Result{
		Timestamp: start,
		Duration:  time.Since(start),
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
