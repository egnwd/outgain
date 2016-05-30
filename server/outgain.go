package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/egnwd/outgain/server/config"
	"github.com/egnwd/outgain/server/engine"
	"github.com/egnwd/outgain/server/routes"

	"github.com/docker/docker/pkg/reexec"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	if reexec.Init() {
		return
	}

	config := config.ParseArgs()

	if config.SandboxMode == "" {
		log.Println("WARNING: sandbox disabled.")
		log.Println("Server is vulnerable to malicious user AIs.")
	} else if config.SandboxMode != "trace" &&
		config.SandboxMode != "kill" &&
		config.SandboxMode != "error" {

		log.Fatal("Invalid sandbox mode: ", config.SandboxMode)
	} else if config.SandboxBin == "" {
		log.Fatal("Sandbox enabled but no sandbox binary given")
	} else if config.SandboxMode == "trace" {
		log.Println("WARNING: sandbox mode \"trace\" is insecure.")
		log.Println("Server is vulnerable to malicious user AIs.")
	}

	engine := engine.NewEngine(config)

	handler := routes.GetHandler(config.StaticDir, engine)
	if config.RedirectPlainHTTP {
		handler = redirectPlainHTTPMiddleware(handler)
	}

	go engine.Run()

	log.Printf("Listening on port %d", config.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", config.Port), handler)
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
