package middleware

import (
	"net/http"
	"context"
	"uuid"
)

type ctxKey int

const (
	requestIDKey ctxKey = iota
)

const (
	requestIdHeader = "X-Request-ID"
)

func RequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate a new request ID
		requestID := r.Header.Get(requestIdHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		w.Header().Set(requestIdHeader, requestID)
		ctx := context.WithValue(r.Context(), requestIDKey, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetRequestIdFromContext(ctx context.Context) string {
	if requestId, ok := ctx.Value(requestIDKey).(string); ok {
		return requestId
	}
	return ""
}