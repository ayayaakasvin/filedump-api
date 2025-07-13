package middlewares

import (
	"net/http"
)

const wildcard = "*"

// CORS middleware
func (m *Middlewares) CorsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if findInSlice(m.allowed_origins, origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
		}

		w.Header().Set("Access-Control-Allow-Methods", m.allowed_methods)
		w.Header().Set("Access-Control-Allow-Headers", m.allowed_headers)
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func findInSlice(slice []string, origin string) bool {
	for _, allowedOrigin := range slice {
		if allowedOrigin == origin || allowedOrigin == wildcard {
			return true
		}
	}

	return false
}