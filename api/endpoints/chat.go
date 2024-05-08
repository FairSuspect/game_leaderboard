package endpoints

import (
	"game_leaderboards/m/v2/api/handlers"

	"github.com/gorilla/mux"
)

func RegisterChatEndpoints(router *mux.Router) {
	router.HandleFunc("/chat", handlers.ChatHandler)
}
