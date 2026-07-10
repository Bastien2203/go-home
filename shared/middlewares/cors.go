package middlewares

import (
	"net/http"
	"strings"

	"github.com/Bastien2203/go-home/shared/config"
)

func CorsMiddleware(allowedOrigins []string, appEnv config.AppEnv) func(http.Handler) http.Handler {
	originSet := make(map[string]bool, len(allowedOrigins))
	for _, o := range allowedOrigins {
		originSet[o] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			if origin != "" {
				allowed := originSet[origin] || (appEnv == config.Dev && IsLocalhostOrigin(origin))
				if allowed {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}
				w.Header().Set("Vary", "Origin")
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Accept, Origin")
			w.Header().Set("Access-Control-Max-Age", "600")

			if strings.HasPrefix(r.URL.Path, "/api") {
				w.Header().Set("Content-Type", "application/json")
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
