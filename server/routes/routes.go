package routes

import (
	"net/http"
	"os"

	"github.com/egnwd/outgain/server/config"
	c "github.com/egnwd/outgain/server/controller"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

//GetHandler returns a mux that maps routes to controller actions
func GetHandler(config *config.Config) http.Handler {
	mux := mux.NewRouter()

	get := mux.Methods(http.MethodGet).Subrouter()
	post := mux.Methods(http.MethodPost).Subrouter()

	get.HandleFunc("/ping", c.PingHandler)

	get.Handle("/images/creature-{colour:[0-9a-fA-F]+}.svg",
		c.SVGSpriteHandler(config.StaticDir))

	get.Handle("/", c.LogInPage(config.StaticDir))

	get.HandleFunc("/login", c.UserLogIn)
	get.Handle("/logout", c.RequireAuth(http.HandlerFunc(c.Logout)))
	get.HandleFunc("/oauthSignInCallback", c.OAuthSignInCallback)
	get.Handle("/currentUser", c.RequireAuth(http.HandlerFunc(c.CurrentUser)))
	get.Handle("/token", c.RequireAuth(http.HandlerFunc(c.UserToken)))

	// Leaderboard
	get.Handle("/peekLeaderboard", c.RequireAuth(http.HandlerFunc(c.LeaderboardPeek)))
	get.Handle("/leaderboard", c.RequireAuth(c.Leaderboard(config.StaticDir)))

	// Leaderboard
	get.Handle("/achievements", c.RequireAuth(c.Achievements(config.StaticDir)))
	get.Handle("/peekAchievements", c.RequireAuth(c.GetAchievements(config.StaticDir)))

	// Lobbies
	get.Handle("/lobbies", c.RequireAuth(c.LobbiesView(config.StaticDir)))
	get.Handle("/peekLobbies", c.RequireAuth(http.HandlerFunc(c.LobbiesPeek)))
	post.Handle("/lobbies/create", c.RequireAuth(c.LobbiesCreate(config)))
	post.Handle("/lobbies/join", c.RequireAuth(http.HandlerFunc(c.LobbiesJoin)))
	post.Handle("/lobbies/leave", c.RequireAuth(http.HandlerFunc(c.LobbiesLeave)))

	// Game View
	lobby := get.PathPrefix("/lobbies/{id:[0-9]+}").Subrouter().StrictSlash(true)
	lobby.Handle("/", c.RequireAuth(c.LobbiesGame(config.StaticDir)))
	lobby.Handle("/users", c.RequireAuth(http.HandlerFunc(c.LobbiesGetUsers)))
	lobby.Handle("/summary", c.RequireAuth(c.LobbiesSummary(config.StaticDir)))
	lobby.Handle("/leaderboard", c.RequireAuth(http.HandlerFunc(c.LobbiesLeaderboard)))
	lobby.Handle("/name", c.RequireAuth(http.HandlerFunc(c.LobbiesName)))

	lobbyPost := post.PathPrefix("/lobbies/{id:[0-9]+}").Subrouter().StrictSlash(true)
	lobbyPost.Handle("/message", c.RequireAuth(http.HandlerFunc(c.LobbiesMessage)))

	get.Handle("/updates/{id:[0-9]+}", c.UpdatesHandler())

	// AI source
	get.Handle("/lobbies/{id:[0-9]+}/ai", c.RequireAuth(c.GetAISource()))
	post.Handle("/lobbies/{id:[0-9]+}/ai", c.RequireAuth(c.PostAISource()))

	// FIXME: Wrap the FileServer in a Handler that hooks w upon writing
	// 404 to the Header
	mux.NotFoundHandler = http.HandlerFunc(c.NotFound)

	get.PathPrefix("/").Handler(
		http.StripPrefix("/", http.FileServer(http.Dir(config.StaticDir))))

	return c.UpdateMaxAge(handlers.LoggingHandler(os.Stdout, mux))
}
