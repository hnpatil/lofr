package lofr

import (
	"context"
	"net/http"
)

func Middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		for k, _ := range r.Header {
			value := r.Header.Get(k)
			ctx = context.WithValue(ctx, k, value)
		}

		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
