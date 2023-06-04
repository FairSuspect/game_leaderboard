package endpoints

import (
	"game_leaderboards/m/v2/api/handlers"
	"game_leaderboards/m/v2/app/repositories"

	"github.com/gorilla/mux"
)

func RegisterLeaderboardEndpoints(router *mux.Router, leaderboardRepository repositories.LeaderboardRepository) {
	lh := handlers.NewLeaderboardHandler(leaderboardRepository)

	router.HandleFunc("/games/{id}/leaderboard", lh.GetGameLeaderboardHandler).Methods("GET")
	router.HandleFunc("/users/{id}/leaderboard", lh.GetUserLeaderboard).Methods("GET")
	router.HandleFunc("/games/{id}/leaderboard", lh.CreateLeaderboardHandler).Methods("POST")
}
