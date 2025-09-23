package api

import (
	"context"
	"net/http"
	"stakeholders-service/internal/models" 

	"github.com/golang-jwt/jwt/v5"
)

// jwtKey je tajni ključ za potpisivanje tokena.
// Pošto se koristi samo u middleware-u, logično je da stoji ovde.
var jwtKey = []byte("super-secret-key")

// AuthMiddleware je funkcija koja presreće zahtev i proverava JWT token.
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Uzmi token iz "Authorization" hedera
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// 2. Ukloni "Bearer " prefiks
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		// 3. Parsiraj i validiraj token
		claims := &models.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// 4. Ako je token validan, dodaj podatke o korisniku u kontekst zahteva
		ctx := r.Context()
		ctx = context.WithValue(ctx, "userID", claims.UserID)
		ctx = context.WithValue(ctx, "userRole", claims.Role)
		r = r.WithContext(ctx)

		// 5. Prosledi zahtev sledećoj funkciji (handleru)
		next.ServeHTTP(w, r)
	}
}