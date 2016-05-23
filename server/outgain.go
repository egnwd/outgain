package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	staticDir := flag.String("static-dir", "client/dist", "")
	redirectPlainHttp := flag.Bool("redirect-plain-http", false, "")
	flag.Parse()

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(*staticDir)))
	mux.HandleFunc("/ping", pingHandler)

	var handler http.Handler = mux
	if *redirectPlainHttp {
		handler = redirectPlainHttpMiddleware(handler)
	}

	http.ListenAndServe(":"+port, handler)
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Pong")
}

func redirectPlainHttpMiddleware(next http.Handler) http.Handler {
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
