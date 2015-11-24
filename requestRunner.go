package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func RunRequest(req *http.Request, res chan Result, nonce bool) {
	ret := runRequest(req, nonce)
	res <- ret
}

func runRequest(req *http.Request, nonce bool) Result {
	if nonce {
		addNonce(req)
	}

	start := time.Now()
	resp, err := client.Do(req)

	result := Result{
		Timestamp: start,
		Duration:  time.Since(start),
		BytesOut:  uint64(req.ContentLength),
	}

	if err != nil {
		result.Error = err.Error()
	} else {
		result.Code = uint16(resp.StatusCode)

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return result
		}

		if result.Code != 200 {
			printError(req, resp, string(body), "Invalid Status Code")
			return result
		}

		if result.Code == 200 {
			result.BytesIn = uint64(len(body))
			result.ReadTime, err = time.ParseDuration(resp.Header.Get("ReadTime"))
			if err != nil {
				printError(req, resp, string(body), "Error parsing read time")
				return result
			}
			result.ReadTimes = resp.Header.Get("ReadTimes")
			result.Server = resp.Header.Get("ServerName")
			result.FanDuration = getDurations(resp)
		}

	}
	return result
}

func getDurations(res *http.Response) []int64 {
	durationsString := res.Header.Get("Durations")
	durationsArray := strings.Split(durationsString, ",")

	var durations = make([]int64, len(durationsArray))
	for i := 0; i < len(durationsArray); i++ {
		durations[i], _ = strconv.ParseInt(durationsArray[i], 10, 64)
	}

	return durations
}

func printError(req *http.Request, res *http.Response, body string, msg string) {
	fmt.Println(strings.Repeat("=", 40))
	if msg != "" {
		fmt.Println("Error: " + msg)
	}
	fmt.Println("Request: " + req.URL.String())
	fmt.Println(strings.Repeat("-", 40))

	fmt.Println("Status: " + res.Status)
	for k, v := range res.Header {
		fmt.Println(k, ":", v)
	}

	fmt.Println(string(body))
	fmt.Println(strings.Repeat("=", 40))
}

func addNonce(req *http.Request) {
	var err error
	rand.Seed(time.Now().UTC().UnixNano())

	queryValues := req.URL.Query()
	apikey := queryValues.Get("apikey")
	nurl := fmt.Sprintf("%s://%s%s?test=%s", req.URL.Scheme, req.URL.Host, req.URL.Path, randomString(100))

	if apikey != "" {
		nurl += "&apikey=" + apikey
	}

	req.URL, err = url.Parse(nurl)
	if err != nil {
		fmt.Println("Error updating url")
		return
	}
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
