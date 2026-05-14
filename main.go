package main

import (
    "fmt"
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
)

func main() {
    fmt.Println("Starting Load Balancer on port 8000...")

    // Parse the string "http://localhost:8081"
	rawStr := "http://localhost:8081"
	myUrl, err := url.Parse(rawStr)
	if err != nil {
		log.Fatal("Something went wrong:", err)
	}

	// Create the reverse proxy
    proxy := httputil.NewSingleHostReverseProxy(myUrl)
	
	// Create the reverse proxy
    proxy := httputil.NewSingleHostReverseProxy(myUrl)
}