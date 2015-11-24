package main

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

type Requestor struct {
	Rate        float64
	Duration    time.Duration
	EndPoint    EndPoint
	ApigeeToken string
	Apigee      ApigeeConfig
	//Statsd StatsdClient
	Config     Config
	Results    Results
	NumResults uint64
	Nonce      bool
}

var t = &http.Transport{
	Dial: func(network, addr string) (net.Conn, error) {
		c, err := net.Dial(network, addr)
		if c == nil {
			fmt.Println("No Connection")
			return c, err
		}

		remote := c.RemoteAddr()
		if remote != nil {
			fmt.Printf("Remote: %s : %s\n", remote.String(), addr)
		}
		return c, err
	},
}
var client = &http.Client{Transport: t}

func (r *Requestor) NewRequest() (http.Request, error) {
	req, err := http.NewRequest("GET", r.Url(), nil)
	if err != nil {
		return *req, err
	}

	if r.ApigeeToken != "" {
		req.Header.Set("Authorization", "Bearer "+r.ApigeeToken)
	}

	return *req, err
}

func NewRequestor(config *Config) (Requestor, error) {
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
		Rate:        config.Rate,
		Duration:    config.Duration,
		EndPoint:    config.Endpoint,
		ApigeeToken: apigeeToken,
		Apigee:      config.Apigee,
		Config:      *config,
		Nonce:       config.Nonce,
	}

	return requestor, nil
}

func (r *Requestor) Url() string {
	if r.ApigeeToken != "" {
		return fmt.Sprintf("%s/%s?apikey=%s", r.Apigee.Apiurl, r.EndPoint.Route, r.Apigee.Apikey)
	} else {
		return fmt.Sprintf("%s://%s:%s/%s", r.EndPoint.Protocol, r.EndPoint.Host, r.EndPoint.Port, r.EndPoint.Route)
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
		//r.Statsd.Timing(r.Results[i].Duration, r.Results[i].ReadTime)
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
			req.Header.Set("Authorization", "Bearer "+token)
		}
	}
}

func runRequests(rate float64, req *http.Request, res chan Result, total uint64, done chan string, nonce bool) {
	fnano := float64(time.Second.Nanoseconds())
	throttle := time.Tick(time.Duration(fnano / rate))
	fmt.Println("Throttle:", time.Duration(1e9/rate))

	for i := 0; uint64(i) < total; i++ {
		<-throttle
		go RunRequest(req, res, nonce)
	}

	done <- "done"
}
