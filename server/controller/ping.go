package controller

import (
	"fmt"
	"net/http"
)

// PingHandler is used for responsiveness checks
func PingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Pong")
}
