package auth

import (
	"net/http"
)

func AuthRequest(next func(http.ResponseWriter, *http.Request, []string)) func(http.ResponseWriter, *http.Request, []string) {
	return func(w http.ResponseWriter, r *http.Request, segments []string) {
		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			w.WriteHeader(http.StatusNetworkAuthenticationRequired)
			w.Write([]byte("Authorization header is required"))

			return
		}
		if !ApiKeyStore.IsValid(authorizationHeader) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))

			return
		}

		next(w, r, segments)
	}
}
