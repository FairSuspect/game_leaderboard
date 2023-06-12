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
	"strconv"
	"unicode/utf8"
)

type UserHandler struct {
	userRepository repository.UserRepository
}

func NewUserHandler(userRepository repository.UserRepository) *UserHandler {
	return &UserHandler{
		userRepository: userRepository,
	}
}

func (uh *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
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
	// Декодируем JSON в структуру User
	var user models.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}
	if user.Name == nil {
		errorMsg := errors.New("Field 'name' is required in model User")
		err := api_errors.DecodeError{Err: errorMsg}
		HandleError(err, w)
		return
	}
	if utf8.RuneCountInString(*user.Name) < 2 {
		errorMsg := errors.New("Name length must at least 2 characters")
		err := api_errors.DecodeError{Err: errorMsg}
		HandleError(err, w)
		return
	}
	createdUser, err := uh.userRepository.CreateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonData, err := json.MarshalIndent(createdUser, "", "	")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		http.Error(w, "Failed to encode user", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)

}

func (uh *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	// Получение параметра id из маршрута
	id, shouldReturn := utils.ParseIdFromRequest(r)
	if shouldReturn {
		return
	}

	user, err := uh.userRepository.GetUserByID(id)
	if err != nil {
		HandleError(err, w)
		return
	}

	// Кодируем срез users в формат JSON
	jsonData, err := serializers.FormatJSON(user)
	if err != nil {
		HandleError(err, w)
		return
	}

	// Устанавливаем заголовок Content-Type на application/json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)

}

func (uh *UserHandler) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	page, err := strconv.Atoi(queryParams.Get("page"))
	if err != nil {
		page = 1
	}

	page_size, err := strconv.Atoi(queryParams.Get("pageSize"))
	if err != nil {
		page_size = 10
	}
	users, err := uh.userRepository.GetAllUsers(page, page_size)
	if err != nil {
		log.Println("Failed to get users:", err)
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		return
	}

	// Кодируем срез users в формат JSON
	jsonData, err := serializers.FormatJSON(users)

	if err != nil {
		log.Println("Failed to format users:", err)
		http.Error(w, "Failed to format users", http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовок Content-Type на application/json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
func (uh *UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Логика обновления пользователя
}

func (uh *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id, shouldReturn := utils.ParseIdFromRequest(r)
	if shouldReturn {
		return
	}
	err := uh.userRepository.DeleteUser(id)
	if err != nil {
		log.Println("Failed to delete user:", err)
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)

	// Логика удаления пользователя
}
