package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// Claims struktura za JWT payload
type Claims struct {
	UserID   int    `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GetJWTSecret čita JWT_SECRET iz environment varijable
func GetJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Fallback za development (ali loguj warning)
		fmt.Println("WARNING: JWT_SECRET not set in environment, using default (UNSAFE for production)")
		return []byte("default-secret-change-this")
	}
	return []byte(secret)
}

// JWTAuthMiddleware proverava validnost JWT tokena i prosleđuje user info dalje
func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Uzmi Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error": "Authorization header required"}`, http.StatusUnauthorized)
			return
		}

		// Ukloni "Bearer " prefix
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			// Ako nije bilo "Bearer " prefiksa
			http.Error(w, `{"error": "Invalid authorization header format"}`, http.StatusUnauthorized)
			return
		}

		// Parsuj i validuj token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Proveri algoritam
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return GetJWTSecret(), nil
		})

		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "Invalid token: %s"}`, err.Error()), http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, `{"error": "Token is not valid"}`, http.StatusUnauthorized)
			return
		}

		// Dodaj user info u context da mikroservisi mogu da ga koriste
		ctx := r.Context()
		ctx = context.WithValue(ctx, "userID", claims.UserID)
		ctx = context.WithValue(ctx, "username", claims.Username)
		ctx = context.WithValue(ctx, "userRole", claims.Role)

		// Prosle��i user info preko headera ka mikroservisima
		r.Header.Set("X-User-ID", fmt.Sprintf("%d", claims.UserID))
		r.Header.Set("X-User-Username", claims.Username)
		r.Header.Set("X-User-Role", claims.Role)

		// Nastavi sa zahtevom
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalJWTMiddleware - ne zahteva token, ali ako postoji validira ga
func OptionalJWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader != "" {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return GetJWTSecret(), nil
			})

			if err == nil && token.Valid {
				// Token je validan, dodaj user info
				ctx := r.Context()
				ctx = context.WithValue(ctx, "userID", claims.UserID)
				ctx = context.WithValue(ctx, "username", claims.Username)
				ctx = context.WithValue(ctx, "userRole", claims.Role)

				r.Header.Set("X-User-ID", fmt.Sprintf("%d", claims.UserID))
				r.Header.Set("X-User-Username", claims.Username)
				r.Header.Set("X-User-Role", claims.Role)

				r = r.WithContext(ctx)
			}
		}

		next.ServeHTTP(w, r)
	})
}
