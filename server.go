package main

    import (
        "fmt"
        "log"
        "net/http"
        "os"
    )

    func main() {
        // We look for an Environment Variable called SERVER_NAME.
        // This is how we will tell our identical containers apart!
        serverName := os.Getenv("SERVER_NAME")
        if serverName == "" {
            serverName = "Unknown Server"
        }

        // When a request comes in, reply with the server's name
        http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
            fmt.Fprintf(w, "Hello from %s\n", serverName)
        })

        log.Printf("Starting %s on port 8080...", serverName)
        log.Fatal(http.ListenAndServe(":8080", nil))
    }