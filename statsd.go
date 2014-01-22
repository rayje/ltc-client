package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

type StatsdClient struct {
	Client     string
	ClientZone string
	Target     string
	TargetZone string
	Type       string
	Conn       net.Conn
}

func NewStatsd(config *Config) StatsdClient {
	statsdConfig := config.Statsd

	if statsdConfig.Host == "" {
		fmt.Println("Statsd Error: Target not set")
		os.Exit(1)
	}
	if statsdConfig.Port == "" {
		fmt.Println("Statsd Error: Port not set")
		os.Exit(1)
	}

	conn, err := net.Dial("udp", statsdConfig.Host+":"+statsdConfig.Port)
	if err != nil {
		panic(err)
	}

	return StatsdClient{
		Client:     config.Client,
		ClientZone: config.Zone,
		Target:     config.Endpoint.Name,
		TargetZone: config.Endpoint.Zone,
		Type:       config.Endpoint.Route,
		Conn:       conn,
	}
}

func (c *StatsdClient) Timing(duration time.Duration, readTime time.Duration) {
	payload := "latency."
	payload += c.Client + "-" + c.ClientZone + "."
	payload += c.Target + "-" + c.TargetZone + "."
	payload += c.Type + ":"
	// payload += strconv.Itoa(int(duration / time.Millisecond))
	payload += strconv.Itoa(int(duration.Nanoseconds()))
	payload += "|ms"

	c.Send(payload)

	readPayload := "latency."
	readPayload += c.Client + "-" + c.ClientZone + "."
	readPayload += c.Target + "-" + c.TargetZone + "."
	readPayload += c.Type + "-readtime:"
	// readPayload += strconv.Itoa(int(duration / time.Millisecond))
	readPayload += strconv.Itoa(int(readTime.Nanoseconds()))
	readPayload += "|ms"

	c.Send(readPayload)
}

func (c *StatsdClient) Send(payload string) {
	c.Conn.Write([]byte(payload))
}
