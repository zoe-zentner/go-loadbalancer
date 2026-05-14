package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

type Backend struct {
	URL   *url.URL
	Alive bool
	mux   sync.RWMutex
}

func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.Alive = alive
}

func (b *Backend) IsAlive() bool {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.Alive
}

// checkHealth tries to open a TCP connection to the backend
func (b *Backend) CheckHealth() bool {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", b.URL.Host, timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

type ServerPool struct {
	mux      sync.Mutex
	Backends []*Backend
	Current  int
}

// Next returns the URL of the next server
func (s *ServerPool) Next() *url.URL {
	s.mux.Lock()
	defer s.mux.Unlock()

	// Loop through all backends to find a healthy one
	for i := 0; i < len(s.Backends); i++ {
		next := (s.Current + i) % len(s.Backends)
		if s.Backends[next].IsAlive() {
			s.Current = (next + 1) % len(s.Backends)
			return s.Backends[next].URL
		}
	}
	return nil
}

// HealthCheck loops through all backends and updates their status
func (s *ServerPool) HealthCheck() {
	for _, b := range s.Backends {
		status := b.CheckHealth()
		b.SetAlive(status)

		msg := "online"
		if !status {
			msg = "offline"
		}
		fmt.Printf("Status Check: %s is %s\n", b.URL.Host, msg)
	}
}

func main() {
	fmt.Println("Starting Load Balancer on port 8000...")

	serverList := []string{
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
	}

	backends := make([]*Backend, 0, len(serverList))

	for _, rawStr := range serverList {
		myUrl, err := url.Parse(rawStr)
		if err != nil {
			log.Fatalf("Invalid backend URL %s: %v", rawStr, err)
		}
		backends = append(backends, &Backend{
			URL:   myUrl,
			Alive: true,
		})
	}

	pool := ServerPool{
		Backends: backends,
		Current:  0,
	}

	fmt.Println("Performing initial health check...")
    pool.HealthCheck()
	
    ticker := time.NewTicker(10 * time.Second)
    go func() {
        // This loop runs every time the ticker sends a value
        for range ticker.C {
            fmt.Println("Starting background health check...")
            pool.HealthCheck()
        }
    }()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		targetUrl := pool.Next()

		// If targetUrl is nil, all backends are down
		if targetUrl == nil {
			http.Error(w, "Service Unavailable: No healthy backends", http.StatusServiceUnavailable)
			return
		}

		fmt.Printf("Directing traffic to: %s\n", targetUrl.Host)
		proxy := httputil.NewSingleHostReverseProxy(targetUrl)
		proxy.ServeHTTP(w, r)
	})

	log.Fatal(http.ListenAndServe(":8000", nil))
}