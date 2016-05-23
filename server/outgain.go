package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	staticDir := flag.String("static-dir", "client/dist", "")

	flag.Parse()

	http.Handle("/", http.FileServer(http.Dir(*staticDir)))
	http.HandleFunc("/ping", pingHandler)

	http.ListenAndServe(":"+port, nil)
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Pong")
}
