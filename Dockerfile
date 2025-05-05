# Use Go 1.23 as the base image
FROM golang:1.24.2-alpine as builder


# Install necessary dependencies
RUN apk --no-cache add ca-certificates tzdata

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files into the container
COPY go.mod go.sum ./

# Download the Go dependencies
RUN go mod download

# Tidy up go.mod and go.sum to fix any issues
RUN go mod tidy

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Start a new stage to minimize the image size
FROM alpine:latest

# Install necessary dependencies for running the app
RUN apk --no-cache add ca-certificates

# Set the working directory inside the container
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Set the entrypoint to run the application
ENTRYPOINT ["./main"]

# Expose the port your app will be listening on
EXPOSE 8080
