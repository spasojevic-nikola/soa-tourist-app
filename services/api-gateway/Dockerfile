# api-gateway/Dockerfile
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Kopiraj go.mod i go.sum
COPY go.mod go.sum ./
RUN go mod download

# Kopiraj sav kod
COPY . .

# Kompajliraj aplikaciju
RUN CGO_ENABLED=0 go build -o /api-gateway main.go

# Finalna, manja slika
FROM alpine:latest AS stage-1

# Instaliraj sertifikate za HTTPS pozive (dobra praksa)
RUN apk add --no-cache ca-certificates

WORKDIR /root/

# Kopiraj kompajlirani binarni fajl
COPY --from=builder /api-gateway .

# Izvršni fajl
CMD ["./api-gateway"]