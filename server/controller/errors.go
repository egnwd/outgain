package controller

import (
	"fmt"
	"net/http"
)

// NotFound is called when a request URI does not exist
func NotFound(w http.ResponseWriter, r *http.Request) {
	fmt.Println("404 Not Found")
	fmt.Fprintf(w, "404 Not found!")
}
