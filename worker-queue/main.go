package main

import (
	"flag"
	"fmt"
	"net/http"
)

var (
	NWorkers = flag.Int("n", 4, "The number of workers to start")
	HTTPAddr = flag.String("http", "127.0.0.1:8000", "Address to listen for HTTP requests on")
)

func init() {
	flag.Parse()
}

func main() {
	fmt.Println("Starting the dispatcher")
	StartDispatcher(*NWorkers)

	fmt.Println("Registering the collector")
	http.HandleFunc("/work", Collector)

	fmt.Println("HTTP server listening on", *HTTPAddr)
	if err := http.ListenAndServe(*HTTPAddr, nil); err != nil {
		fmt.Println(err.Error())
	}
}
