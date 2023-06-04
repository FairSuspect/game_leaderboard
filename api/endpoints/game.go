package endpoints

import (
	"game_leaderboards/m/v2/api/handlers"
	"game_leaderboards/m/v2/app/repositories"

	"github.com/gorilla/mux"
)

func RegisterGameEndpoints(router *mux.Router, gameRepository repositories.GameRepository) {
	uh := handlers.NewGameHandler(gameRepository)
	router.HandleFunc("/games", uh.CreateGameHandler).Methods("POST")
	router.HandleFunc("/games", uh.GetGamesHandler).Methods("GET")
	router.HandleFunc("/games/{id}", uh.GetGameHandler).Methods("GET")
	router.HandleFunc("/games/{id}", uh.UpdateGameHandler).Methods("PUT")
	router.HandleFunc("/games/{id}", uh.DeleteGameHandler).Methods("DELETE")
}
