package main

import (
	"fmt"
	"net/http"
	"time"
	"flag"
)



func main() {
	r := flag.String("r", "small", "The route to call on the server")
	h := flag.String("h", "localhost", "The host of the server")
	p := flag.String("p", "8080", "The port of the host server")
	flag.Parse()

	url := "http://"+ *h +":" + *p + "/" + *r

	makeRequest(url)
}