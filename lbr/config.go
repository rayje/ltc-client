package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

type ApigeeConfig struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Apikey   string `json:"apikey"`
	Apiurl   string `json:"apiurl"`
}

type EndPoint struct {
	Route string `json:"route"`
	Host  string `json:"host"`
	Port  string `json:"port"`
	Protocol string `json:protocol"`
}

type Config struct {
	Apigee         ApigeeConfig  `json:"apigee"`
	UseApigee      bool          `json:"useapigee"`
	EndPoint       EndPoint      `json:"endpoint"`
}

func getConfig() Config {
	route := flag.String("route", "small", "The route to call on the server (small|med|large|xlarge)")
	host := flag.String("host", "localhost", "The host of the server")
	port := flag.String("port", "80", "The port of the host server")
	https := flag.Bool("https", false, "Use https as the transfer protocol.")

	apigee := flag.Bool("apigee", false, "Use an apigee request")
	configFile := flag.String("config", "config.json", "Location of config file")

	flag.Parse()

	var config Config

	if *apigee {
		file, err := ioutil.ReadFile(*configFile)
		if err != nil {
			fmt.Printf("File error: %v\n", err)
			os.Exit(1)
		}

		err = json.Unmarshal(file, &config)
		if err != nil {
			fmt.Printf("JSON error: %v\n", err)
			os.Exit(1)
		}
	}

	config.setEndpoint(*route, *host, *port, *https)
	config.UseApigee = *apigee

	return config
}

func (c *Config) setEndpoint(route string, host string, port string, https bool) {
	var emptyEndPoint = EndPoint{}

	if c.EndPoint == emptyEndPoint {
		endpoint := EndPoint{
			Route: route,
			Host:  host,
			Port:  port,
		}
		c.EndPoint = endpoint
	}

	if https {
		c.EndPoint.Protocol = "https"
	} else {
		c.EndPoint.Protocol = "http"
	}
}