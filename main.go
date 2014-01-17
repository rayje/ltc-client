package main

import (
	"fmt"
	"time"
	"flag"
)

func main() {
	route := flag.String("route", "small", "The route to call on the server")
	host := flag.String("host", "localhost", "The host of the server")
	port := flag.String("port", "8080", "The port of the host server")
	rate := flag.Uint64("rate", 1, "Requests per second")
	duration := flag.Duration("duration", 1*time.Second, "Duration of the test")
	flag.Parse()

	url := "http://"+ *host +":" + *port + "/" + *route

	res, err := makeRequest(url, *rate, *duration)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res)
}