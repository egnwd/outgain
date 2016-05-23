package routes

import (
	"net/http"

	"github.com/egnwd/outgain/server/controller"
	"github.com/gorilla/mux"
)

//GetHandler returns a mux that mapps routes to controller actions
func GetHandler(static string) http.Handler {
	mux := mux.NewRouter()

	get := mux.Methods("GET").Subrouter()
	// post := mux.Methods("POST").Subrouter()

	get.Handle("/", http.FileServer(http.Dir(static)))
	get.HandleFunc("/ping", controller.PingHandler)
	get.HandleFunc("/signin", controller.UserSignIn)

	return mux
}
