# Use the official lightweight Go image
FROM golang:1.22-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy your server.go file from your laptop into the container
COPY server.go .

# Initialize a Go module and build the application into an executable called 'server'
RUN go mod init dummy-server && go build -o server server.go

# Document that this container listens on port 8080 internally
EXPOSE 8080

# The command to run when the container starts
CMD ["./server"]