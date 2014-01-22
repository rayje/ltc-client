ltc-client
==========

A Go client for latency testing

## Usage:
	$ ltc-client -h
	Usage of ltc-client:
	  -client="localhost": The name of the client server
	  -config="config.json": Location of config file
	  -duration=1s: Duration of the test
	  -host="localhost": The host of the server
	  -port="8080": The port of the host server
	  -rate=1: Requests per second
	  -route="small": The route to call on the server