package middleware

import (
	"net/http"
	"strings"
)

type CORSConfig struct {
	AllowedOrigins     []string
	AllowedMethods     []string
	AllowedHeaders     []string
	AllowedCredentials bool
}

func DefaultCorsConfig() CORSConfig {
	return CORSConfig{
		AllowedOrigins:     []string{"*"},
		AllowedMethods:     []string{"POST", "PUT", "PATCH", "DELETE", "OPTIONS", "GET"},
		AllowedHeaders:     []string{"Content-Type", "Authorization", "Accept"},
		AllowedCredentials: false,
	}
}

func ProductionCORSConfig(origins []string) CORSConfig {
	return CORSConfig{
		AllowedOrigins:     origins,
		AllowedMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowedCredentials: true,
	}
}

func CORS(config CORSConfig) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			allowed := false

			for _, allowedOrigin := range config.AllowedOrigins {
				if allowedOrigin == "*" || allowedOrigin == origin {
					allowed = true
					break
				}
			}
			if allowed {
				if origin != "" {
					w.Header().Set("Access-Control-Allowed-Origin", origin)
				} else if len(config.AllowedOrigins) == 1 && config.AllowedOrigins[0] == "*" {
					w.Header().Set("Access-Control-Allow-Origin", "*")
				}
			}

			w.Header().Set("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ", "))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ", "))

			if config.AllowedCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
			next(w, r)
		}
	}
}
