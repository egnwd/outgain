package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"routes"
	"strings"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	staticDir := flag.String("static-dir", "client/dist", "")
	redirectPlainHTTP := flag.Bool("redirect-plain-http", false, "")
	flag.Parse()

	handler := routes.GetHandler(*staticDir)
	if *redirectPlainHTTP {
		handler = redirectPlainHTTPMiddleware(handler)
	}

	http.ListenAndServe(":"+port, handler)
}

func redirectPlainHTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-forwarded-proto") != "https" {
			hostname := strings.Split(r.Host, ":")[0]
			redirectTo := fmt.Sprintf("https://%s%s", hostname, r.URL.String())
			http.Redirect(w, r, redirectTo, http.StatusFound)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
