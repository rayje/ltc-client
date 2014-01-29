package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
    "math/rand"
    "net/url"

)

type Requestor struct {
	Rate     float64
	Duration time.Duration
	EndPoint EndPoint
	ApigeeToken string
	Apigee ApigeeConfig
	Statsd StatsdClient
	Config Config
	Results Results
	NumResults uint64
	Nonce bool
}

type Result struct {
	Code      uint16
	Timestamp time.Time
	Duration  time.Duration
	BytesOut  uint64
	BytesIn   uint64
	Error     string
	ReadTime  time.Duration
}

type Results []Result

var client = &http.Client{}

func (r *Requestor) NewRequest() (http.Request, error) {
	req, err := http.NewRequest("GET", r.Url(), nil)
	if err != nil {
		return *req, err
	}

	if r.ApigeeToken != "" {
		req.Header.Set("Authorization", "Bearer " + r.ApigeeToken)
	}

	return *req, err
}

func NewRequestor(config *Config, statsd StatsdClient) (Requestor, error) {
	var apigeeToken string
	var requestor Requestor
	var err error

	if config.UseApigee {
		apigeeToken, err = getApigeeToken(config)
		if err != nil {
			return requestor, err
		}
	}

	requestor = Requestor{
		Rate:     config.Rate,
		Duration: config.Duration,
		EndPoint: config.Endpoint,
		ApigeeToken: apigeeToken,
		Apigee: config.Apigee,
		Statsd: statsd,
		Config: *config,
		Nonce: config.Nonce,
	}

	return requestor, nil
}

func (r *Requestor) Url() string {
	if r.ApigeeToken != "" {
		return fmt.Sprintf("%s/%s?apikey=%s", r.Apigee.Apiurl, r.EndPoint.Route, r.Apigee.Apikey)
	} else {
		return fmt.Sprintf("http://%s:%s/%s", r.EndPoint.Host, r.EndPoint.Port, r.EndPoint.Route)
	}
}

func (r *Requestor) MakeRequest() (Results, error) {
	total := uint64(r.Rate * float64(r.Duration.Seconds()))
	fmt.Println("Total Requests:", total)

	res := make(chan Result, total)
	done := make(chan string)
	r.Results = make(Results, total)

	req, err := r.NewRequest()
	if err != nil {
		return nil, err
	}

	if r.ApigeeToken != "" {
		go tokenRefresh(&r.Config, &req, done)
	}

	go runRequests(r.Rate, &req, res, total, done, r.Nonce)

	for i := 0; i < cap(res); i++ {
		r.Results[i] = <-res
		r.NumResults += 1
		r.Statsd.Timing(r.Results[i].Duration, r.Results[i].ReadTime)
	}
	close(res)

	return r.Results, nil
}

func (r *Requestor) GetResults() Results {
	return r.Results[:r.NumResults]
}

func tokenRefresh(config *Config, req *http.Request, done chan string) {
	refresh := time.Tick(time.Hour - (5 * time.Minute))

	for {
		select {
		case <-done:
		    return
		case <-refresh:
			token, err := getApigeeToken(config)
			if err != nil {
				fmt.Println(err)
				return
			}
			req.Header.Set("Authorization", "Bearer " + token)
		}
	}
}

func runRequests(rate float64, req *http.Request, res chan Result, total uint64, done chan string, nonce bool) {
	throttle := time.Tick(time.Duration(1e9 / rate))
	fmt.Println("Throttle:", time.Duration(1e9 / rate))

	for i := 0; uint64(i) < total; i++ {
		<-throttle
		go runRequest(req, res, nonce)
	}

	done <- "done"
}

func runRequest(req *http.Request, res chan Result, nonce bool) {
	if nonce {
		var err error
	    rand.Seed( time.Now().UTC().UnixNano())
	    req.URL, err = url.Parse(req.URL.String() + "?test=" + randomString(100))

	    if err != nil {
	    	fmt.Println("Error updating url")
	    	return
	    }
	}

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
			fmt.Println(err)
		} else {
			if result.Code < 200 || result.Code >= 300 {
				fmt.Println("======================================")
				fmt.Println("Status: " + r.Status)
				for k, v := range r.Header {
					fmt.Println(k, ":", v)
				}
				fmt.Println(string(body))
				fmt.Println("======================================")
			} else {
				result.BytesIn = uint64(len(body))
				result.ReadTime, err = time.ParseDuration(r.Header.Get("ReadTime"))
				if err != nil {
					fmt.Println("======================================")
					fmt.Println("Error parsing read time")
					fmt.Println("--------------------------------------")
					fmt.Println("Status: " + r.Status)
					for k, v := range r.Header {
						fmt.Println(k, ":", v)
					}
					fmt.Println(string(body))
					fmt.Println("======================================")
				}
			}
		}
	}

	res <- result
}

func randomString (l int ) string {
    bytes := make([]byte, l)
    for i:=0 ; i<l ; i++ {
        bytes[i] = byte(randInt(65,90))
    }
    return string(bytes)
}

func randInt(min int , max int) int {
    return min + rand.Intn(max-min)
}
