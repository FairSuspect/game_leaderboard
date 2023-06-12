package main

import (
	"database/sql"
	"fmt"

	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"game_leaderboards/m/v2/api/endpoints"
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
	gameRepository := repository.NewGameRepository(db)
	leaderboardRepository := repository.NewLeaderboardRepository(db)

	// Инициализация маршрутизатора

	// Регистрация обработчиков
	endpoints.RegisterUserEndpoints(r, userRepository)
	endpoints.RegisterGameEndpoints(r, gameRepository)
	endpoints.RegisterLeaderboardEndpoints(r, leaderboardRepository)

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
