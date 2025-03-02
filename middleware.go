package lofr

import (
	"context"
	"net/http"
)

const Headers = "http_headers"

func Middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), Headers, r.Header)))
	})
}
