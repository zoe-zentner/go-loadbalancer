# Layer 7 Go Load Balancer

A Layer 7 reverse proxy and load balancer written in Go.

It routes incoming HTTP traffic across a pool of Dockerized backend servers, ensuring high availability even if individual nodes fail.

## Core Features

- **Round Robin Routing:** Distributes incoming HTTP traffic evenly across a pool of healthy backend servers.
- **Concurrency Safe:** Utilizes `sync.Mutex` and `sync.RWMutex` to prevent race conditions during simultaneous traffic spikes and backend state changes.
- **Active Fault Tolerance:** A background Goroutine continuously pings backends via TCP. Offline servers are dynamically removed from the routing pool and re-added automatically upon recovery.
- **Dynamic Configuration:** Reads target URLs, ports, and timeout settings from a central JSON file, requiring no code recompilation to alter network behavior.

## Useful Docker Commands

The backend servers are managed via Docker Compose. Here are the commands to control the environment:

- **Start the environment (in background):**

  ```bash
  docker compose up --build -d
  ```

- **View live logs of all backend servers:**

  ```bash
  docker compose logs -f
  ```

- **Stop a specific server (useful for testing failure):**

  ```bash
  docker stop load-balancer-server_b-1
  ```

  _(You can replace this with your exact container name if it differs.)_

- **Restart a stopped server (useful for testing recovery):**

  ```bash
  docker start load-balancer-server_b-1
  ```

- **Shut down the entire environment:**
  ```bash
  docker compose down
  ```

## Getting Started

### 1. Start the Backend Pool

Spin up the three identical Go web servers for testing.

```bash
docker compose up -d
```

### 2. Configure the Load Balancer

Create or edit the `config.json` file to control the load balancer's behavior. You can adjust the health check interval and add or remove backends here:

```json
{
  "port": ":8000",
  "backends": [
    "http://localhost:8081",
    "http://localhost:8082",
    "http://localhost:8083"
  ],
  "timeout_seconds": 2,
  "health_check_interval_seconds": 10
}
```

### 3. Run the Load Balancer

Start the application. You can specify a custom config file path using the `-config` flag.

```bash
go run main.go -config=config.json
```

### 4. Test Fault Tolerance (Chaos Testing)

To see the dependable systems logic in action, send a continuous stream of requests to the load balancer:

```bash
# On Linux/macOS/WSL
while true; do curl http://localhost:8000; sleep 1; done
```

While the loop is running, open a new terminal tab and stop one of your Docker containers (e.g., `docker stop load-balancer-server_b-1`). Within 10 seconds, the load balancer's background checker will detect the TCP timeout, mark the node as dead, and dynamically route all new traffic to the surviving nodes without dropping any user requests.
