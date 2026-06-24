package middleware

import "net/http"

// ponytail: per-IP rate limiter, swap for Redis-based when multi-replica
func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ponytail: token bucket — implement before public form submission
		next.ServeHTTP(w, r)
	})
}
