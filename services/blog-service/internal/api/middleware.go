package api

import (
"context"
"net/http"
"strconv"
)

func AuthMiddleware(next http.Handler) http.Handler {
return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
userIDStr := r.Header.Get("X-User-Id")
if userIDStr == "" {
http.Error(w, "Missing X-User-Id header", http.StatusUnauthorized)
return
}
userID, err := strconv.ParseUint(userIDStr, 10, 32)
if err != nil {
http.Error(w, "Invalid X-User-Id header", http.StatusBadRequest)
return
}

ctx := context.WithValue(r.Context(), "userID", uint(userID))
next.ServeHTTP(w, r.WithContext(ctx))
})
}
