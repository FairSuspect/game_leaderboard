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
	"net/http"
	"strconv"
)

type LeaderboardHandler struct {
	leaderboardRepository repository.LeaderboardRepository
}

func NewLeaderboardHandler(leaderboardRepository repository.LeaderboardRepository) *LeaderboardHandler {
	// Создание объекта с начальными значениями (выделение памяти)
	leaderboardHandler := new(LeaderboardHandler)
	leaderboardHandler.leaderboardRepository = leaderboardRepository
	return leaderboardHandler
}

func (uh *LeaderboardHandler) CreateLeaderboardHandler(w http.ResponseWriter, r *http.Request) {
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
	// Декодируем JSON в структуру Leaderboard
	var leaderboardEntry models.LeaderboardEntry
	err = json.Unmarshal(body, &leaderboardEntry)
	if err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}
	if leaderboardEntry.Score < 0 {
		errorMsg := errors.New("Field 'score' must be positive")
		err := api_errors.DecodeError{Err: errorMsg}
		HandleError(err, w)
		return
	}

	err = uh.leaderboardRepository.AddRecord(leaderboardEntry)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

}

func (uh *LeaderboardHandler) GetGameLeaderboardHandler(w http.ResponseWriter, r *http.Request) {
	// Получение параметра id из маршрута
	id, shouldReturn := utils.ParseIdFromRequest(r)
	if shouldReturn {
		return
	}
	queryParams := r.URL.Query()
	page, err := strconv.Atoi(queryParams.Get("page"))
	if err != nil {
		page = 1
	}

	page_size, err := strconv.Atoi(queryParams.Get("pageSize"))
	if err != nil {
		page_size = 10
	}

	leaderboard, err := uh.leaderboardRepository.GetGameLeaderboard(id, page, page_size)
	if err != nil {
		HandleError(err, w)
		return
	}

	// Кодируем срез leaderboards в формат JSON
	jsonData, err := serializers.FormatJSON(leaderboard)
	if err != nil {
		HandleError(err, w)
		return
	}

	// Устанавливаем заголовок Content-Type на application/json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)

}
func (uh *LeaderboardHandler) GetUserLeaderboard(w http.ResponseWriter, r *http.Request) {
	// Получение параметра id из маршрута
	id, shouldReturn := utils.ParseIdFromRequest(r)
	if shouldReturn {
		return
	}

	leaderboard, err := uh.leaderboardRepository.GetUserLeaderboard(id)
	if err != nil {
		HandleError(err, w)
		return
	}

	// Кодируем срез leaderboards в формат JSON
	jsonData, err := serializers.FormatJSON(leaderboard)
	if err != nil {
		HandleError(err, w)
		return
	}

	// Устанавливаем заголовок Content-Type на application/json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)

}

func (uh *LeaderboardHandler) UpdateLeaderboardHandler(w http.ResponseWriter, r *http.Request) {
	// Логика обновления пользователя
}
