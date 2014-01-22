package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type StatsdConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type EndPoint struct {
	Route string `json:"route"`
	Host  string `json:"host"`
	Port  string `json:"port"`
}

type Config struct {
	Statsd   StatsdConfig  `json:"statsd"`
	Endpoint EndPoint      `json:"endpoint"`
	Rate     uint64        `json:"rate"`
	Duration time.Duration `json:"duration"`
	Client   string        `json:"client"`
}

func getConfig() Config {
	route := flag.String("route", "small", "The route to call on the server")
	host := flag.String("host", "localhost", "The host of the server")
	port := flag.String("port", "8080", "The port of the host server")
	rate := flag.Uint64("rate", 1, "Requests per second")
	client := flag.String("client", "localhost", "The name of the client server")
	duration := flag.Duration("duration", 1*time.Second, "Duration of the test")
	configFile := flag.String("config", "config.json", "Location of config file")
	flag.Parse()

	file, err := ioutil.ReadFile(*configFile)
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}

	var config Config
	json.Unmarshal(file, &config)

	config.setEndpoint(*route, *host, *port)
	config.setRateDuration(*rate, *duration)
	config.setClient(*client)

	return config
}

func (c *Config) setEndpoint(route string, host string, port string) {
	var emptyEndPoint = EndPoint{}

	if c.Endpoint == emptyEndPoint {
		endpoint := EndPoint{
			Route: route,
			Host:  host,
			Port:  port,
		}
		c.Endpoint = endpoint
	}
}

func (c *Config) setRateDuration(rate uint64, duration time.Duration) {
	if c.Rate == 0 {
		c.Rate = rate
	}

	if c.Duration == 0 {
		c.Duration = duration
	}
}

func (c *Config) setClient(client string) {
	if c.Client == "" {
		c.Client = client
	}
}