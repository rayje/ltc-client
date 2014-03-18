package main

import (
	"fmt"
	"net/http"
	"os"
	"net"
	"time"
	"bytes"
	"flag"
)

var timeout = time.Duration(5 * time.Second)

func DialMethod(network, addr string) (net.Conn, error) {
	c, err := net.DialTimeout(network, addr, timeout)
	if c == nil {
		fmt.Printf("No Connection ")
		return c, err
	}

	remote := c.RemoteAddr()
	if remote != nil {
		fmt.Printf("Remote: %s : %s ",remote.String(), addr)
	} else {
		fmt.Printf("RemoteAddr not set ")
	}
	return c, err
}

func getUrl() string {
	route := flag.String("route", "med", "The route to call on the server (small|med|large|xlarge)")
	host := flag.String("host", "localhost", "The host of the server")
	port := flag.String("port", "80", "The port of the host server")
	flag.Parse()

	var buffer bytes.Buffer
	buffer.WriteString("http://")
	buffer.WriteString(*host)
	if *port != "80" {
		buffer.WriteString(":")
		buffer.WriteString(*port)
	}
	buffer.WriteString("/")
	buffer.WriteString(*route)

	return buffer.String()
}

func main() {
	var t =  &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: DialMethod,
	}
	var client = &http.Client{Transport: t}
	var url = getUrl()

	req, e := http.NewRequest("GET", url, nil)
	if e != nil {
		fmt.Println("Error creating request")
		fmt.Println(e)
		os.Exit(1)
	}

	start := time.Now()
	resp, err := client.Do(req)
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("Error retrieving response\n")
		fmt.Println(err)
		return
	}

	fmt.Printf("%d %v\n", resp.StatusCode, duration)
}
