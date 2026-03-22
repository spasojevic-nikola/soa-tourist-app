package api

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"tour-service/internal/models"
)

// UserContextMiddleware injects user claims (if present) into the request context so handlers
// can consistently read them without duplicating header parsing logic.
func UserContextMiddleware(next http.Handler) http.Handler {
	secret := []byte(os.Getenv("JWT_SECRET"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		hasUserContext := false

		if userIDStr := r.Header.Get("X-User-ID"); userIDStr != "" {
			if parsed, err := strconv.ParseUint(userIDStr, 10, 64); err == nil {
				ctx = context.WithValue(ctx, "userID", uint(parsed))
				hasUserContext = true
			}
		}

		if username := r.Header.Get("X-User-Username"); username != "" {
			ctx = context.WithValue(ctx, "username", username)
			hasUserContext = true
		}

		if role := r.Header.Get("X-User-Role"); role != "" {
			ctx = context.WithValue(ctx, "userRole", role)
			hasUserContext = true
		}

		if !hasUserContext {
			if claims, ok := parseClaimsFromAuth(r.Header.Get("Authorization"), secret); ok {
				ctx = context.WithValue(ctx, "userID", claims.UserID)
				ctx = context.WithValue(ctx, "username", claims.Username)
				ctx = context.WithValue(ctx, "userRole", claims.Role)

				r.Header.Set("X-User-ID", strconv.FormatUint(uint64(claims.UserID), 10))
				r.Header.Set("X-User-Username", claims.Username)
				r.Header.Set("X-User-Role", claims.Role)
			}
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func parseClaimsFromAuth(authHeader string, secret []byte) (*models.Claims, bool) {
	if len(secret) == 0 {
		return nil, false
	}

	if authHeader == "" || !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
		return nil, false
	}

	tokenString := strings.TrimSpace(authHeader[len("Bearer "):])
	claims := &models.Claims{}
	parsedToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenUnverifiable
		}
		return secret, nil
	})

	if err != nil || !parsedToken.Valid {
		return nil, false
	}

	return claims, true
}

// getUserIDFromHeader izvlači User ID iz X-User-ID headera (postavljenog od API Gateway-a)
func GetUserIDFromHeader(r *http.Request) (int, error) {
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		return 0, nil
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

// GetUserRoleFromHeader izvlači User Role iz X-User-Role headera
func GetUserRoleFromHeader(r *http.Request) string {
	return r.Header.Get("X-User-Role")
}

// getUsernameFromHeader izvlači Username iz X-User-Username headera
func GetUsernameFromHeader(r *http.Request) string {
	return r.Header.Get("X-User-Username")
}
