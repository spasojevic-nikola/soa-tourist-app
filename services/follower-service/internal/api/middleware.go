package api

import (
	"context"
	"net/http"
	"strconv"
)

// Definišemo custom tip za kljuc u kontekstu da izbegnemo kolizije
type contextKey string
const userKey contextKey = "userID"

// AuthMiddleware cita hedere koje postavlja API Gateway i stavlja userID u kontekst
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Citamo heder X-User-ID koji ce postaviti API Gateway
		userIDStr := r.Header.Get("X-User-ID")
		if userIDStr == "" {
			http.Error(w, "Unauthorized: Missing user information from gateway", http.StatusUnauthorized)
			return
		}

		// Konvertujemo ID iz stringa u uint
		userID, err := strconv.ParseUint(userIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Unauthorized: Invalid user ID format", http.StatusUnauthorized)
			return
		}

		// Stavljamo userID u kontekst zahteva
		ctx := context.WithValue(r.Context(), userKey, uint(userID))
		
		// Prosledjujemo zahtev sa novim kontekstom sledecem handleru
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}