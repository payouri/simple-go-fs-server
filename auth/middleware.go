package auth

import (
	"net/http"
)

func AuthRequest(next func(http.ResponseWriter, *http.Request, []string)) func(http.ResponseWriter, *http.Request, []string) {
	return func(w http.ResponseWriter, r *http.Request, segments []string) {
		if r.Header.Get("Authorization") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next(w, r, segments)
	}
}
