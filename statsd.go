package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

type StatsdClient struct {
	Client string
	Host   string
	conn   net.Conn
}

func NewStatsd(config *Config) StatsdClient {
	statsdConfig := config.Statsd

	if statsdConfig.Host == "" {
		fmt.Println("Statsd Error: Host not set")
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

	return StatsdClient{config.Client, config.Endpoint.Host, conn}
}

func (c *StatsdClient) Timing(duration time.Duration) {
	payload := "latency." + c.Client + "." + c.Host
	payload += ".timing:"
	// payload += strconv.Itoa(int(duration/time.Millisecond))
	payload += strconv.Itoa(int(duration.Nanoseconds()))
	payload += "|ms"

	c.Send(payload)
}

func (c *StatsdClient) Send(payload string) {
	fmt.Println("Sending: " + payload)
	c.conn.Write([]byte(payload))
}
