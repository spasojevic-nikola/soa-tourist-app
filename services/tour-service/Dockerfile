# Build stage - KORISTIMO NOVU VERZIJU GO-A
FROM golang:1.25-alpine AS builder

# Potrebno za go mod download
RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build aplikacije
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Kopiranje binarnog fajla iz buildera
COPY --from=builder /app/main .

# Port koji aplikacija sluša unutar kontejnera
EXPOSE 8080

# Pokretanje aplikacije
CMD ["./main"]