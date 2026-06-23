package api

import (
	"net/http"
	"strings"
)

func AuthMiddleware(token string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/healthz" {
			next.ServeHTTP(w, r)
			return
		}

		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") || strings.TrimPrefix(header, "Bearer ") != token {
			writeJSON(w, http.StatusUnauthorized, map[string]bool{"ok": false})
			return
		}

		next.ServeHTTP(w, r)
	})
}
