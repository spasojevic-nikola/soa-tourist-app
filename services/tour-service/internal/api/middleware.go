package api

import (
	"context"
	"net/http"
	"tour-service/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("super-secret-key")

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if len(tokenString) < 8 || tokenString[:7] != "Bearer " {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}
		tokenString = tokenString[7:]
		claims := &models.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		ctx = context.WithValue(ctx, "userRole", claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func AuthorOrAdminAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userRole, ok := r.Context().Value("userRole").(string)
		if !ok || (userRole != "author" && userRole != "administrator") {
			http.Error(w, "Forbidden: Authors or administrators only", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	}
}