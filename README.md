# Load Balancer Project

This project is a Layer 7 Load Balancer built in Go, designed to route traffic across multiple backend Docker containers.

## Useful Docker Commands

- **Start the environment (in background):**
  `docker compose up --build -d`
- **View live logs of all servers:**
  `docker compose logs -f`
- **Shut down the environment:**
  `docker compose down`
- **Shut down and wipe out volumes/networks (clean slate):**
  `docker compose down -v`

## Testing the Backend

To verify the dummy servers are running, hit their individual mapped ports:

- `curl http://localhost:8081` (Server A)
- `curl http://localhost:8082` (Server B)
- `curl http://localhost:8083` (Server C)

## Running the Load Balancer

The load balancer implements a **Round Robin** algorithm with **Active Health Checking**.

- **Concurrency:** Uses Mutexes and RWMutexes to prevent race conditions during state changes.
- **Fault Tolerance:** A background Goroutine pings backend servers via TCP every 10 seconds. If a server goes offline, it is dynamically removed from the routing pool. When it recovers, it is added back.

**To Test Fault Tolerance:**

1. Start the environment: `docker compose up -d`
2. Start the proxy: `go run main.go`
3. Send continuous requests to `http://localhost:8000`.
4. Kill a backend container (e.g., `docker stop server_b`) and observe the traffic automatically route around the failure.
