# # 1. BUILD STAGE (Faza izgradnje)
# --------------------------
    FROM golang:1.25-alpine AS builder

    # Instalacija git-a je neophodna za "go mod download" i za dohvat privatnih/zasebnih repozitorijuma
    RUN apk add --no-cache git
    
    # Postavljanje radnog direktorijuma unutar kontejnera
    WORKDIR /app
    
    # Kopiranje go.mod i go.sum fajlova.
    # Ovo omogućava Docker-u da kešira preuzimanje modula ako se kod promeni, ali ne i zavisnosti.
    COPY go.mod go.sum ./
    
    # Preuzimanje zavisnosti (modula) i popravljanje/provera.
    # Zahvaljujući go mod tidy, ovaj korak je sada brz.
    RUN go mod download
    RUN go mod tidy
    
    # Kopiranje ostatka izvornog koda
    COPY . .
    
    # Izgradnja binarne aplikacije.
    # CGO_ENABLED=0 osigurava statički binarni fajl (bez eksternih zavisnosti OS-a).
    # ./cmd/api je putanja do tvog main.go fajla.
    RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api
    
    # --------------------------
    # # 2. FINAL STAGE (Konačna faza)
    # --------------------------
    FROM alpine:latest
    
    # Instalacija CA sertifikata (potrebno za HTTPS/TLS komunikaciju, npr. sa drugim servisima)
    RUN apk --no-cache add ca-certificates
    
    # Postavljanje radnog direktorijuma za izvršavanje
    WORKDIR /root/
    
    # Kopiranje binarnog fajla "main" iz faze izgradnje
    COPY --from=builder /app/main .
    
    # Ekspozovanje porta koji koristi tvoja Go aplikacija (8081)
    EXPOSE 8081
    
    # Komanda za pokretanje binarnog fajla kada se kontejner startuje
    CMD ["./main"]