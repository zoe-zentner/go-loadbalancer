package main

import (
    "fmt"
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
	"sync"
)

// ServerPool holds our list of backends and tracks the current one
type ServerPool struct {
	mux sync.Mutex
	Backends []*url.URL
	Current int
}

// Next returns the next server in the list using Round Robin
func (s *ServerPool) Next() *url.URL {
	s.mux.Lock()
    defer s.mux.Unlock() 
	target := s.Backends[s.Current]
	s.Current = (s.Current + 1) % len(s.Backends)
	return target
}

func main() {
    fmt.Println("Starting Load Balancer on port 8000...")

	// List of backend servers to balance across
	serverList := []string{
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
	}

	// Slice to hold the parsed URL objects
	backends := make([]*url.URL, 0, len(serverList))

	// Convert raw strings into URL objects
	for _, rawStr := range serverList {
        myUrl, err := url.Parse(rawStr)
        if err != nil {
            log.Fatalf("Invalid backend URL %s: %v", rawStr, err)
        }
        backends = append(backends, myUrl)
    }

	// Initialize the server pool
	pool := ServerPool{
		Backends: backends,
		Current: 0,
	}

	// Handler that runs for every curl request
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		target := pool.Next()
		fmt.Printf("Directing traffic to: %s\n", target.Host)
		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.ServeHTTP(w, r)
	})

	// Start the load balancer server
    log.Fatal(http.ListenAndServe(":8000", nil))
}