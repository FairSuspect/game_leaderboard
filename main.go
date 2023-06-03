package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"game_leaderboards/m/v2/api/endpoints"
	"game_leaderboards/m/v2/app/models"
	repository "game_leaderboards/m/v2/app/repositories"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "fair"
	password = "fair"
	dbname   = "game"
)

var db *sql.DB

type UserHandler struct {
	userRepository repository.UserRepository
}

func main() {
	// Подключение к базе данных
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := mux.NewRouter()
	userRepository := repository.NewUserRepository(db)

	// Инициализация маршрутизатора

	// Регистрация обработчиков
	endpoints.RegisterUserEndpoints(r, userRepository)

	// r.HandleFunc("/users", userHandler.usersRouteHandler).Methods("GET", "POST")
	// r.HandleFunc("/user/{id}", userHandler.userRouteHandler).Methods("GET", "PUT")
	// r.HandleFunc("/leaderboard", userHandler.userRouteHandler).Methods("GET", "POST", "PUT")

	// http.HandleFunc("/create_user", createUserHandler)
	// http.HandleFunc("/create_leaderboard", createLeaderboardHandler)
	// http.HandleFunc("/update_leaderboard", updateLeaderboardHandler)
	// http.HandleFunc("/list_leaderboard", listLeaderboardHandler)
	log.Println("Starting web server on port 8888")
	routes := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, _ := route.GetPathTemplate()
		fmt.Println("Path:", path)
		return nil
	})

	if routes != nil {
		fmt.Println(routes)
	}

	// Запуск веб-сервера на порту 8080
	log.Fatal(http.ListenAndServe(":8888", r))

}

func (uh *UserHandler) usersRouteHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		uh.getUsers(w, r)
	case http.MethodPost:
		uh.createUser(w, r)
	default:
		// Обработка других методов
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
func (uh *UserHandler) userRouteHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		uh.getUser(w, r)
	case http.MethodPut:
		uh.editUser(w, r)
	default:
		// Обработка других методов
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (uh *UserHandler) getUsers(w http.ResponseWriter, r *http.Request) {
	users, err := uh.userRepository.GetAllUsers()
	if err != nil {
		log.Println("Failed to get users:", err)
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		return
	}

	// Кодируем срез users в формат JSON
	jsonData, err := json.MarshalIndent(users, "", "  ")

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

func (uh *UserHandler) getUser(w http.ResponseWriter, r *http.Request) {
	// Получение параметра id из маршрута
	vars := mux.Vars(r)
	idRaw := vars["id"]

	id, err := strconv.Atoi(idRaw)
	if err != nil {
		fmt.Println("Ошибка преобразования строки в int:", err)
		return
	}

	user, err := uh.userRepository.GetUserByID(id)

	// Кодируем срез users в формат JSON
	jsonData, err := json.MarshalIndent(user, "", "	")
	if err != nil {
		log.Println("Failed to format user:", err)
		http.Error(w, "Failed to format user", http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовок Content-Type на application/json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)

}
func (uh *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {

	// Считываем тело запроса
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// Декодируем JSON в структуру User
	var user models.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}
	createdUser, err := uh.userRepository.CreateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonData, err := json.MarshalIndent(createdUser, "", "	")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		http.Error(w, "Failed to encode user", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)

}
func (uh *UserHandler) editUser(w http.ResponseWriter, r *http.Request) {

	// Считываем тело запроса
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// Декодируем JSON в структуру User
	var user models.User
	err = json.Unmarshal(body, &user)
	vars := mux.Vars(r)
	idRaw := vars["id"]

	id, err := strconv.Atoi(idRaw)
	user.ID = id
	if err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}
	newUser, err := uh.userRepository.EditUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		http.Error(w, "Failed to create user", http.StatusBadRequest)
		return
	}
	jsonData, err := json.MarshalIndent(newUser, "", "	")

	w.Write(jsonData)
	w.WriteHeader(http.StatusOK)

}

func createLeaderboardHandler(w http.ResponseWriter, r *http.Request) {
	// Получение параметров "user_id" и "score" из URL-запроса
	userID := r.URL.Query().Get("user_id")
	score := r.URL.Query().Get("score")

	// Вставка новой записи лидера в таблицу "leaderboard"
	_, err := db.Exec("INSERT INTO leaderboard (user_id, score) VALUES ($1, $2)", userID, score)
	if err != nil {
		log.Println("Failed to create leaderboard:", err)
		http.Error(w, "Failed to create leaderboard", http.StatusInternalServerError)
		return
	}

	// Отправка успешного ответа
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Leaderboard created successfully")
}

func updateLeaderboardHandler(w http.ResponseWriter, r *http.Request) {
	// Получение параметров "user_id" и "score" из URL-запроса
	userID := r.URL.Query().Get("user_id")
	score := r.URL.Query().Get("score")

	// Обновление записи лидера в таблице "leaderboard"
	result, err := db.Exec("UPDATE leaderboard SET score = $1 WHERE user_id = $2", score, userID)
	if err != nil {
		log.Println("Failed to update leaderboard:", err)
		http.Error(w, "Failed to update leaderboard", http.StatusInternalServerError)
		return
	}
	// Проверка, была ли обновлена запись
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Failed to get rows affected:", err)
		http.Error(w, "Failed to update leaderboard", http.StatusInternalServerError)
		return
	}

	// Проверка, была ли обновлена хотя бы одна запись
	if rowsAffected == 0 {
		log.Println("No leaderboard record updated")
		http.Error(w, "No leaderboard record updated", http.StatusNotFound)
		return
	}

	// Отправка успешного ответа
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Leaderboard updated successfully")
}

func listLeaderboardHandler(w http.ResponseWriter, r *http.Request) {
	// Получение параметров "page" и "limit" из URL-запроса
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	// Установка значений по умолчанию, если параметры не заданы
	if page == "" {
		page = "1"
	}
	if limit == "" {
		limit = "10"
	}
	// Выполнение SQL-запроса для постраничной выборки лидеров
	rows, err := db.Query("SELECT users.username, leaderboard.score FROM leaderboard INNER JOIN users ON leaderboard.user_id = users.user_id ORDER BY leaderboard.score DESC LIMIT $1 OFFSET $2", limit, page)
	if err != nil {
		log.Println("Failed to list leaderboard:", err)
		http.Error(w, "Failed to list leaderboard", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Перебор результатов и формирование ответа
	var leaderboardList []string
	for rows.Next() {
		var username string
		var score int
		err := rows.Scan(&username, &score)
		if err != nil {
			log.Println("Failed to scan leaderboard row:", err)
			http.Error(w, "Failed to list leaderboard", http.StatusInternalServerError)
			return
		}
		leaderboardList = append(leaderboardList, fmt.Sprintf("Username: %s, Score: %d", username, score))
	}

	// Проверка наличия ошибок при переборе результатов
	err = rows.Err()
	if err != nil {
		log.Println("Failed to iterate leaderboard rows:", err)
		http.Error(w, "Failed to list leaderboard", http.StatusInternalServerError)
		return
	}

	// Отправка успешного ответа
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Leaderboard:")
	for _, entry := range leaderboardList {
		fmt.Fprintln(w, entry)
	}
}
