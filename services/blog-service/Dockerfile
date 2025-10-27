# Build stage
FROM golang:1.25-alpine AS builder

# Install git for go mod download. This is a good practice.
RUN apk add --no-cache git

WORKDIR /app

# Copy go.mod and go.sum first for caching
COPY go.mod go.sum ./

# Download dependencies (update go.sum if needed)
RUN go mod download
RUN go mod tidy

# Copy all source code
COPY . .

# Ensure go.sum is complete after copying
RUN go mod tidy

# Build the application with CGO disabled for a static binary.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Final stage - a minimal alpine image
FROM alpine:latest

# Install CA certificates to enable secure connections (HTTPS/TLS)
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Expose the port used by the application
EXPOSE 8081

# Command to run the binary
CMD ["./main"]