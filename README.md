ltc-client
==========

A Go client for latency testing

This client will make HTTP requests to a remote host, capturing the round trip time and calculating measurements to determine variablility.

The measurements captured are:

* 25 Percentile
* 75 Percentile
* 95 Percentile
* 99 Percentile
* 999 Percentile
* Mean
* Standard Deviation
* Min
* Max
* Inter-Quartile range
	

The following is an example of the reports generated from the captured round trip times.

	$ ltc-client -duration 3h -host 10.10.10.10 -port 80 -route med -rate 0.1 -rint 0
	Total Requests: 1080
	Throttle: 10s
	==============================
	Results - final
	------------------------------
	Requests: 1080
	Latencies:
		Total:     5.159692077s
		     	   Client	    Fan
				------------	------------
		0.25:	3.862988ms		0
		0.75:	4.164414ms		0
		0.95:	5.466224ms		0
		0.99:	17.764905ms		0
		0.999:	269.887298ms	0
	
		mean:	4.777492ms		0
		std:	9.204613ms		0
		min:	3.467133ms		0
		max:	289.146555ms	0
		iqr:	301.426us		0
	==============================

Along with this report, the capture times are dumped to a file named ```results.txt```. If the ```rint``` flag is not set to ```0```, the dumped times are sorted.

## Usage:

	$ ltc-client -h
	Usage of ltc-client:
	  -apigee=false: Use an apigee request
	  -config="config.json": Location of config file
	  -duration=1s: Duration of the test
	  -fan=false: Report fan in final results
	  -host="localhost": The host of the server
	  -https=false: Use https as the transfer protocol.
	  -nonce=false: Use a nonce for each request
	  -port="8080": The port of the host server
	  -rate=1: Requests per second
	  -rint=15m0s: Interval to print reports
	  -route="small": The route to call on the server (small|med|large|xlarge)

### apigee (default: false)
A boolean flag to indicate whether a request through apigee is required. 

If this flag is set to true, then a config file that contains apigee information is required. The following
is an example of the apigee content required in the config file.

	{
		"apigee": {
			"email": "...",
			"password": "...",
			"apikey": "...",
			"apiurl": "..."
		}
	}

### config (default: config.json)
The location of the config file. The default location is the location where the command is run, where the default config file name is ```config.json```

**Note**: The config file is only required if ```apigee``` is set to true.

### duration (default: 1s)
The amount of time that the test should run. The time formats used are based on the GoLang's ```time.Duration``` The following format are acceptable:

	1h - 1 hour
	1m - 1 minute
	1s - 1 second
	1ms - 1 millisecond
	1us - 1 microsecond
	1ns - 1 nanosecond
	
### fan
In the case where the HTTP request is made to an ```ltc-mid``` server, then the ```-fan``` flag will indicate that a fan report should be calculated.

### host (default: localhost)
The host that should be used for HTTP requests.

### https (default: false)
A boolean flag to indicate whether the protocal used for the HTTP request should be HTTPS.

### nonce (default: false)
A boolean flag to indicate whether a unique value should be appended to the HTTP request to help prevent caching.

### port (default: 8080)
The port to be used when making HTTP requests.

### rate (default: 1s)
The number of requests to make per second. This can also be indicated as a fraction to run requests at a rate greater than 1 second.

	0.1 - 1 request every 10 seconds
	10  - 10 requests per second

### rint (default: 15m0s)
The ```rint``` flag is used to indicate how often an updated report should be printed. This flag will force intermediate reported to be printed to the console. **Note**: Set this flag to ```0``` to disable intermediate reports.

### route (default: small)
The route to call on the host. This will resolve to ```http://localhost:8080/small```.
