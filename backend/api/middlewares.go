package api

import (
	"net/http"
	"strings"
)

// middleware: strip trailing slashes
func stripSlashes(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimRight(r.URL.Path, "/")
		}

		next.ServeHTTP(w, r)
	})
}

// middleware: restrict methods
func restrictMethod(method string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}

		next.ServeHTTP(w, r)
	})
}
