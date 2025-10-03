package api

import (
	"context"
	"net/http"
	"strconv"
)

// AuthMiddleware osigurava da je korisnik autentifikovan.
// Za sada, MOCK: čita X-User-ID header postavljen od strane API Gateway-a.
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDHeader := r.Header.Get("X-User-ID")
		
		if userIDHeader == "" {
			http.Error(w, "Authorization header or X-User-ID required", http.StatusUnauthorized)
			return
		}
		
		userID, err := strconv.ParseUint(userIDHeader, 10, 64)
		if err != nil {
			http.Error(w, "Invalid User ID format", http.StatusUnauthorized)
			return
		}
		
		// Postavljanje UserID-ja u kontekst
		ctx := r.Context()
		ctx = context.WithValue(ctx, "userID", uint(userID))
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}
}

// helper funkcija za dobijanje UserID iz konteksta
func GetUserID(r *http.Request) uint {
	// U produkcijskom kodu bi se proveravalo 'ok'
	userID, _ := r.Context().Value("userID").(uint)
	return userID
}