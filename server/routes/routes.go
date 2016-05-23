package routes

import (
	"fmt"
	"net/http"
	"strings"
)

func getHandler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(*staticDir)))
	mux.HandleFunc("/ping", pingHandler)

	var handler http.Handler = mux
	if *redirectPlainHttp {
		handler = redirectPlainHttpMiddleware(handler)
	}

	return handler
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
