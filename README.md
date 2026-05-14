# Load Balancer Project - FYP Prep

This project is a Layer 7 Load Balancer built in Go, designed to route traffic across multiple backend Docker containers.

## Useful Docker Commands

*   **Start the environment (in background):** 
    `docker compose up --build -d`
*   **View live logs of all servers:** 
    `docker compose logs -f`
*   **Shut down the environment:** 
    `docker compose down`
*   **Shut down and wipe out volumes/networks (clean slate):** 
    `docker compose down -v`

## Testing the Backend
To verify the dummy servers are running, hit their individual mapped ports:
*   `curl http://localhost:8081` (Server A)
*   `curl http://localhost:8082` (Server B)
*   `curl http://localhost:8083` (Server C)

## Running the Load Balancer

Currently, the load balancer acts as a simple reverse proxy.

1. Ensure the backend Docker environment is running.
2. Start the proxy locally:
   `go run main.go`
3. Test the proxy:
   `curl http://localhost:8000`

*Note: Right now, the proxy statically routes all traffic to Server A on port 8081.*