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
	Name  string `json:"name"`
	Zone  string `json:"zone"`
	Type  string `json:"zone"`
}

type Config struct {
	Statsd    StatsdConfig  `json:"statsd"`
	Apigee    ApigeeConfig  `json:"apigee"`
	Endpoint  EndPoint      `json:"endpoint"`
	Rate      float64       `json:"rate"`
	Duration  time.Duration `json:"duration"`
	Client    string        `json:"client"`
	Zone      string        `json:"zone"`
	UseApigee bool          `json:"useapigee"`
	ReportInterval time.Duration`json:"reportInterval"`
}

func getConfig() Config {
	route := flag.String("route", "small", "The route to call on the server (small|med|large|xlarge)")
	host := flag.String("host", "localhost", "The host of the server")
	port := flag.String("port", "8080", "The port of the host server")

	rate := flag.Float64("rate", 1.0, "Requests per second")
	duration := flag.Duration("duration", 1*time.Second, "Duration of the test")

	client := flag.String("client", "localhost", "The name of the client server")
	clientzone := flag.String("clientzone", "us-east-1b", "The name of the client server")
	target := flag.String("target", "localhost", "The name of the target (for graphite)")
	targetzone := flag.String("targetzone", "us-east-1b", "The name of the aws zone (for graphite)")

	reportInterval := flag.Duration("rint", 15*time.Minute, "Interval to print reports")
	apigee := flag.Bool("apigee", false, "Use an apigee request")
	configFile := flag.String("config", "config.json", "Location of config file")
	flag.Parse()

	file, err := ioutil.ReadFile(*configFile)
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}

	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		fmt.Printf("JSON error: %v\n", err)
		os.Exit(1)
	}

	config.setEndpoint(*route, *host, *port, *target, *targetzone)
	config.setRateDuration(*rate, *duration)
	config.setClient(*client, *clientzone)
	config.setReportIntercal(*reportInterval)
	config.UseApigee = *apigee

	return config
}

func (c *Config) setEndpoint(route string, host string, port string, name string, zone string) {
	var emptyEndPoint = EndPoint{}

	if c.Endpoint == emptyEndPoint {
		endpoint := EndPoint{
			Route: route,
			Host:  host,
			Port:  port,
			Name:  name,
			Zone:  zone,
		}
		c.Endpoint = endpoint
	}
}

func (c *Config) setRateDuration(rate float64, duration time.Duration) {
	if c.Rate == 0 {
		c.Rate = rate
	}

	if c.Duration == 0 {
		c.Duration = duration
	}
}

func (c *Config) setClient(client string, zone string) {
	if c.Client == "" {
		c.Client = client
	}

	if c.Zone == "" {
		c.Zone = zone
	}
}

func (c *Config) setReportIntercal(interval time.Duration) {
	if (c.ReportInterval == 0) {
		c.ReportInterval = interval
	}
}
