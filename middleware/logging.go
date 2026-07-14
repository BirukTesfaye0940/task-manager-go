package middleware

import (
	"log"
	"net/http"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			reqID := GetRequestID(r.Context())
			if reqID != "" {
				log.Printf("[%s] %s %s", reqID, r.Method, r.URL.Path)
			} else {
				log.Printf("%s %s", r.Method, r.URL.Path)
			}

			next.ServeHTTP(w, r)
		},
	)
}
