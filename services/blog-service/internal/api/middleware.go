package api

import (
	"context"
	"net/http"
	"strings"

	"blog-service/internal/models" 

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("super-tajni-kljuc-koji-niko-ne-zna-12345")

// AuthMiddleware proverava validnost JWT tokena.
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims := &models.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "userID", claims.UserID)
		ctx = context.WithValue(ctx, "userRole", claims.Role)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}
}