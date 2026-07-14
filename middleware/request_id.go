package middleware

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
)

type contextKey string

const RequestIDKey contextKey = "request_id"

var reqIDCounter uint64

// RequestID injects a unique request ID into the request context.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := atomic.AddUint64(&reqIDCounter, 1)
		reqID := fmt.Sprintf("req-%d", id)

		// Set the X-Request-ID response header for client transparency
		w.Header().Set("X-Request-ID", reqID)

		ctx := context.WithValue(r.Context(), RequestIDKey, reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetRequestID retrieves the request ID from the context, or empty string if not present.
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return ""
}
