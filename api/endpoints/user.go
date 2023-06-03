package endpoints

import (
	"game_leaderboards/m/v2/api/handlers"
	"game_leaderboards/m/v2/app/repositories"

	"github.com/gorilla/mux"
)

func RegisterUserEndpoints(router *mux.Router, userRepository repositories.UserRepository) {
	uh := handlers.NewUserHandler(userRepository)
	router.HandleFunc("/users", uh.CreateUserHandler).Methods("POST")
	router.HandleFunc("/users", uh.GetUsersHandler).Methods("GET")
	router.HandleFunc("/users/{id}", uh.GetUserHandler).Methods("GET")
	router.HandleFunc("/users/{id}", uh.UpdateUserHandler).Methods("PUT")
	router.HandleFunc("/users/{id}", uh.DeleteUserHandler).Methods("DELETE")
}
