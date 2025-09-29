# Build faza
FROM golang:1.25-alpine AS builder

# Instalacija potrebnih alata
RUN apk add --no-cache git

# Postavljanje radnog direktorijuma
WORKDIR /app

# Kopiranje modula i preuzimanje zavisnosti
COPY go.mod go.sum ./
RUN go mod download

# Kopiranje celokupnog izvornog koda
COPY . .

# Izgradnja (build) aplikacije
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Finalna faza
FROM alpine:latest

# Kopiranje sertifikata za HTTPS pozive
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Kopiranje izgrađene aplikacije iz build faze
COPY --from=builder /app/main .

# Izlaganje porta
EXPOSE 8080

# Komanda za pokretanje aplikacije
CMD ["./main"]