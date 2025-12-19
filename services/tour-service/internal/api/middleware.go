package api

import (
	"net/http"
	"strconv"
)

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
