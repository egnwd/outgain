package routes

import (
	"net/http"
	"os"

	c "github.com/egnwd/outgain/server/controller"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

//GetHandler returns a mux that maps routes to controller actions
func GetHandler(static string) http.Handler {
	mux := mux.NewRouter()

	get := mux.Methods(http.MethodGet).Subrouter()
	post := mux.Methods(http.MethodPost).Subrouter()

	get.HandleFunc("/ping", c.PingHandler)

	get.Handle("/images/creature-{id:[0-9a-fA-F]+}.png",
		c.SpriteHandler(static))

	get.Handle("/", c.LogInPage(static))

	get.HandleFunc("/login", c.UserLogIn)
	get.Handle("/logout", c.RequireAuth(http.HandlerFunc(c.Logout)))
	get.HandleFunc("/oauthSignInCallback", c.OAuthSignInCallback)
	get.HandleFunc("/currentUser", c.CurrentUser)

	// Lobbies
	get.Handle("/lobbies", c.RequireAuth(c.LobbiesView(static)))
	get.Handle("/peekLobbies", c.RequireAuth(http.HandlerFunc(c.LobbiesPeek)))
	get.Handle("/lobbies/{id:[0-9]+}/users", c.RequireAuth(http.HandlerFunc(c.LobbiesGetUsers)))
	post.Handle("/lobbies/join", c.RequireAuth(http.HandlerFunc(c.LobbiesJoin)))

	// Game View
	get.Handle("/lobbies/{id:[0-9]+}", c.RequireAuth(c.LobbiesGame(static)))
	get.Handle("/updates/{id:[0-9]+}", c.UpdatesHandler())
	get.HandleFunc("/leave", c.Leave)

	// FIXME: Wrap the FileServer in a Handler that hooks w upon writing
	// 404 to the Header
	mux.NotFoundHandler = http.HandlerFunc(c.NotFound)

	get.PathPrefix("/").Handler(
		http.StripPrefix("/", http.FileServer(http.Dir(static))))

	return c.UpdateMaxAge(handlers.LoggingHandler(os.Stdout, mux))
}
