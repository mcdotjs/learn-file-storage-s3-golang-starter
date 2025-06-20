package main

import (
	"fmt"
	"net/http"
)

func noCacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// w.Header().Set("Cache-Control", "max-age=3600")
		w.Header().Set("Cache-Control", "no-store")
		fmt.Println("cache")
		next.ServeHTTP(w, r)
	})
}
