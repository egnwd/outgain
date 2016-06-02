package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/egnwd/outgain/server/config"
	"github.com/egnwd/outgain/server/database"
	"github.com/egnwd/outgain/server/lobby"
	"github.com/egnwd/outgain/server/routes"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	config := config.ParseArgs()

	db, err := database.OpenDb()
	_ = db
	if err != nil {
		log.Fatal(err)
	}

	handler := routes.GetHandler(config.StaticDir)
	if config.RedirectPlainHTTP {
		handler = redirectPlainHTTPMiddleware(handler)
	}

	lobby.GenerateOneLobby(config)

	log.Printf("Listening on port %d", config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), handler))
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
