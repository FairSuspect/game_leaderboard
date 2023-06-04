package handlers

import (
	"encoding/json"
	"errors"
	api_errors "game_leaderboards/m/v2/api/errors"
	serializers "game_leaderboards/m/v2/api/serializers"
	utils "game_leaderboards/m/v2/api/utils"
	"game_leaderboards/m/v2/app/models"
	repository "game_leaderboards/m/v2/app/repositories"
	"io/ioutil"
	"log"
	"net/http"
	"unicode/utf8"
)

type GameHandler struct {
	gameRepository repository.GameRepository
}

func NewGameHandler(gameRepository repository.GameRepository) *GameHandler {
	return &GameHandler{
		gameRepository: gameRepository,
	}
}

func (uh *GameHandler) CreateGameHandler(w http.ResponseWriter, r *http.Request) {
	// Считываем тело запроса
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		HandleError(err, w)
		return
	}

	// Проверяем пустое тело запроса
	if len(body) == 0 {
		http.Error(w, "Request body is empty", http.StatusBadRequest)
		return
	}
	// Декодируем JSON в структуру Game
	var game models.Game
	err = json.Unmarshal(body, &game)
	if err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}
	if game.Name == nil {
		errorMsg := errors.New("Field 'name' is required in model Game")
		err := api_errors.DecodeError{Err: errorMsg}
		HandleError(err, w)
		return
	}
	if utf8.RuneCountInString(*game.Name) < 2 {
		errorMsg := errors.New("Name length must at least 2 characters")
		err := api_errors.DecodeError{Err: errorMsg}
		HandleError(err, w)
		return
	}
	createdGame, err := uh.gameRepository.CreateGame(game)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonData, err := json.MarshalIndent(createdGame, "", "	")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		http.Error(w, "Failed to encode game", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)

}

func (uh *GameHandler) GetGameHandler(w http.ResponseWriter, r *http.Request) {
	// Получение параметра id из маршрута
	id, shouldReturn := utils.ParseIdFromRequest(r)
	if shouldReturn {
		return
	}

	game, err := uh.gameRepository.GetGameByID(id)
	if err != nil {
		HandleError(err, w)
		return
	}

	// Кодируем срез games в формат JSON
	jsonData, err := serializers.FormatJSON(game)
	if err != nil {
		HandleError(err, w)
		return
	}

	// Устанавливаем заголовок Content-Type на application/json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)

}

func (uh *GameHandler) GetGamesHandler(w http.ResponseWriter, r *http.Request) {
	games, err := uh.gameRepository.GetAllGames()
	if err != nil {
		log.Println("Failed to get games:", err)
		http.Error(w, "Failed to get games", http.StatusInternalServerError)
		return
	}

	// Кодируем срез games в формат JSON
	jsonData, err := serializers.FormatJSON(games)

	if err != nil {
		log.Println("Failed to format games:", err)
		http.Error(w, "Failed to format games", http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовок Content-Type на application/json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
func (uh *GameHandler) UpdateGameHandler(w http.ResponseWriter, r *http.Request) {
	// Логика обновления пользователя
}

func (uh *GameHandler) DeleteGameHandler(w http.ResponseWriter, r *http.Request) {
	id, shouldReturn := utils.ParseIdFromRequest(r)
	if shouldReturn {
		return
	}
	err := uh.gameRepository.DeleteGame(id)
	if err != nil {
		log.Println("Failed to delete game:", err)
		http.Error(w, "Failed to delete game", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)

	// Логика удаления пользователя
}
