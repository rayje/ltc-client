package main

import (
	"fmt"
	"net/http"
	"os"
	"net"
	"crypto/tls"
	"time"
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

	c.SetDeadline(time.Now().Add(timeout));
	return c, err
}

func getUrl(config *Config) (string, error) {
	var apigeeToken string
	var err error

	if config.UseApigee {
		apigeeToken, err = getApigeeToken(config)
		if err != nil {
			return "", err
		}
	}

	if apigeeToken != "" {
		return fmt.Sprintf("%s/%s?apikey=%s",
			config.Apigee.Apiurl, config.EndPoint.Route, config.Apigee.Apikey), nil
	} else {
		return fmt.Sprintf("%s://%s:%s/%s",
			config.EndPoint.Protocol, config.EndPoint.Host,
			config.EndPoint.Port, config.EndPoint.Route), nil
	}
}

func getTransport(config *Config) http.RoundTripper {
	var t http.RoundTripper

	if config.UseApigee || config.EndPoint.Protocol == "https" {
		t = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial: DialMethod,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
 	} else {
 		t = &http.Transport{Dial: DialMethod}
 	}

	return t
}

func main() {
	config := getConfig()

	var t = getTransport(&config);
	var client = &http.Client{Transport: t}
	url, err := getUrl(&config)
	if err != nil {
		fmt.Println("Error building URL")
		fmt.Println(err)
		os.Exit(1)
	}

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
