package main

import (
	"encoding/json"
	"flag"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync"
	"time"
)

// Config maps directly to the config.json file
type Config struct {
	Port                string   `json:"port"`
	Backends            []string `json:"backends"`
	TimeoutSecs         int      `json:"timeout_seconds"`
	HealthCheckInterval int      `json:"health_check_interval_seconds"`
}

type Backend struct {
	URL     *url.URL
	Alive   bool
	mux     sync.RWMutex
	Timeout time.Duration
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

func (b *Backend) CheckHealth() bool {
	conn, err := net.DialTimeout("tcp", b.URL.Host, b.Timeout)
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

func (s *ServerPool) Next() *url.URL {
	s.mux.Lock()
	defer s.mux.Unlock()

	for i := 0; i < len(s.Backends); i++ {
		next := (s.Current + i) % len(s.Backends)
		if s.Backends[next].IsAlive() {
			s.Current = (next + 1) % len(s.Backends)
			return s.Backends[next].URL
		}
	}
	return nil
}

func (s *ServerPool) HealthCheck() {
	for _, b := range s.Backends {
		status := b.CheckHealth()
		b.SetAlive(status)

		msg := "online"
		if !status {
			msg = "offline"
		}
		log.Printf("Health Check: %s is %s\n", b.URL.Host, msg)
	}
}

// loadConfig reads a JSON file and converts it into the Config struct
func loadConfig(filepath string) (*Config, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func main() {
	// Setup Command Line Flags
	configPath := flag.String("config", "config.json", "Path to the configuration file")
	flag.Parse()

	// Load Configuration
	cfg, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v\n", err)
	}

	log.Printf("Starting Load Balancer on port %s...\n", cfg.Port)

	// Initialize Backend Pool
	backends := make([]*Backend, 0, len(cfg.Backends))
	timeout := time.Duration(cfg.TimeoutSecs) * time.Second

	for _, rawStr := range cfg.Backends {
		myUrl, err := url.Parse(rawStr)
		if err != nil {
			log.Fatalf("Invalid backend URL %s: %v", rawStr, err)
		}
		backends = append(backends, &Backend{
			URL:     myUrl,
			Alive:   true,
			Timeout: timeout,
		})
	}

	pool := ServerPool{
		Backends: backends,
		Current:  0,
	}

	// Start Health Checker
	log.Println("Performing initial health check...")
	pool.HealthCheck()

	interval := time.Duration(cfg.HealthCheckInterval) * time.Second
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			pool.HealthCheck()
		}
	}()

	// Start Reverse Proxy Server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		targetUrl := pool.Next()

		if targetUrl == nil {
			http.Error(w, "Service Unavailable: No healthy backends", http.StatusServiceUnavailable)
			return
		}

		// Using standard log instead of fmt for cleaner output
		log.Printf("Routing %s request to: %s\n", r.Method, targetUrl.Host)
		proxy := httputil.NewSingleHostReverseProxy(targetUrl)
		proxy.ServeHTTP(w, r)
	})

	log.Fatal(http.ListenAndServe(cfg.Port, nil))
}