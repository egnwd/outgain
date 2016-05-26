package controller

import (
	"fmt"
	"net/http"
)

// Leave temporarily logs the user out - this will change in the future
func Leave(w http.ResponseWriter, r *http.Request) {
	u := fmt.Sprintf("http://%s/logout", r.Host)

	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}
