package middleware

import "net/http"

// JWT verifikasi Supabase — placeholder
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ponytail: validate Supabase JWT from Authorization header
		next.ServeHTTP(w, r)
	})
}
