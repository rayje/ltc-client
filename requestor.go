package main

var client = &http.Client{}

func makeRequest(url string, rate uint64, duration time.Duration) {
	total := rate * uint64(duration.Seconds())
	res := make(chan Result, total)

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Client", "Header")

	startTime := time.Now()
	resp, err := client.Do(req)
	endTime := time.Now()

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Route:", url)
	fmt.Println("==========================")
	fmt.Println("RTT:", endTime.Sub(startTime))

	fmt.Println("\nHeaders")
	fmt.Println("--------------------------")
	for k, v := range resp.Header {
		fmt.Println(k, v)
	}
}