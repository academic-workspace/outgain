package routes

import (
	"net/http"
	"os"

	"github.com/egnwd/outgain/server/controller"
	"github.com/egnwd/outgain/server/engine"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

//GetHandler returns a mux that maps routes to controller actions
func GetHandler(static string, engine *engine.Engine) http.Handler {
	mux := mux.NewRouter()

	get := mux.Methods(http.MethodGet).Subrouter()
	post := mux.Methods(http.MethodPost).Subrouter()

	get.HandleFunc("/ping", controller.PingHandler)

	get.HandleFunc("/images/creature-{id:[0-9a-fA-F]+}.png",
		controller.SpriteHandler(static))

	get.Handle("/", controller.LogInPage(static))

	get.HandleFunc("/login", controller.UserLogIn)
	get.HandleFunc("/logout", controller.Logout)
	get.HandleFunc("/oauthSignInCallback", controller.OAuthSignInCallback)
	get.HandleFunc("/currentUser", controller.CurrentUser)

	// Lobbies
	get.HandleFunc("/lobbies", controller.LobbiesView)
	post.HandleFunc("/lobbies/join", controller.LobbiesJoin)

	// Game View
	get.HandleFunc("/lobbies/{id:[0-9]+}", controller.LobbiesGame)
	get.Handle("/updates/{id:[0-9]+}", controller.UpdatesHandler(engine))
	get.HandleFunc("/leave", controller.Leave)

	// FIXME: Wrap the FileServer in a Handler that hooks w upon writing
	// 404 to the Header
	mux.NotFoundHandler = http.HandlerFunc(controller.NotFound)

	get.PathPrefix("/").Handler(
		http.StripPrefix("/", http.FileServer(http.Dir(static))))

	return controller.UpdateMaxAge(handlers.LoggingHandler(os.Stdout, mux))
}
