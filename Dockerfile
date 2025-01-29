# Use the official Golang image to create a build artifact.
# Use a specific Go version as needed.
FROM golang:1.23-alpine AS builder
# Set the Current Working Directory inside the container
WORKDIR /app
# Copy go mod and sum files
COPY go.mod go.sum ./
# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download
# Copy the source from the current directory to the Working Directory inside the container
COPY . .
# Build the Go app for Linux OS and amd64 architecture
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o myapp .
# Start a new stage from scratch
FROM alpine:latest
# Set up certificates
RUN apk --no-cache add ca-certificates
# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/myapp /app/myapp
# Expose port 8080 to the outside world
EXPOSE 8080
# Command to run the executable
CMD ["/app/myapp"]
