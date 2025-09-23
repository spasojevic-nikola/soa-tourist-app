# Build stage
FROM golang:1.21-alpine AS builder

# Install git for go mod download. This is a good practice.
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod and go sum files first.
# This improves caching and rebuild speed.
COPY go.mod go.sum ./

# Download and tidy up dependencies.
# The go mod tidy command will automatically add or remove
# any modules required by the source code.
RUN go mod download
RUN go mod tidy

# Copy the rest of the source code.
# This will happen only if go.mod or go.sum changes.
COPY . .

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