package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	models "game_leaderboards/m/v2/app/models"

	api_errors "game_leaderboards/m/v2/api/errors"

	"github.com/lib/pq"
)

type userRepository struct {
	db *sql.DB
}

// UserRepository представляет интерфейс для работы с пользователями в базе данных.
type UserRepository interface {
	GetAllUsers(page int, page_size int) (*[]models.User, error)
	GetUserByID(id int) (*models.User, error)
	CreateUser(user models.User) (*models.User, error)
	EditUser(user models.User) (*models.User, error)
	DeleteUser(id int) error
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (ur *userRepository) CreateUser(user models.User) (*models.User, error) {
	rows := ur.db.QueryRow("INSERT INTO users (name) VALUES ($1) RETURNING ID", user.Name)

	var id int
	err := rows.Scan(&id)
	if err != nil {

		// Проверяем, является ли ошибка ошибкой нарушения уникального ограничения
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Code == "23505" { // код "23505" соответствует ошибке нарушения уникального ограничения
			// Формируем свое сообщение об ошибке
			errMsg := fmt.Sprintf("Username '%s' already exists", *user.Name)
			return nil, errors.New(errMsg)
		}

		// Возвращаем другие ошибки без изменений
		return nil, err

	}
	user.ID = int(id)
	return &user, nil
}
func (ur *userRepository) EditUser(user models.User) (*models.User, error) {
	_, err := ur.db.Exec("UPDATE users SET name = $1 WHERE id = $2", user.Name, user.ID)
	if err != nil {

		// Возвращаем другие ошибки без изменений
		return nil, err

	}
	return &user, err
}

func (ur *userRepository) GetUserByID(id int) (*models.User, error) {
	var user models.User
	err := ur.db.QueryRow("SELECT id, name FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name)
	if err != nil {
		log.Println("Failed to get user by ID:", err)
		return nil, api_errors.NotFoundError{Id: id}
	}
	return &user, nil
}
func (ur *userRepository) DeleteUser(id int) error {
	_, err := ur.db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		log.Println("Failed to get user by ID:", err)
		return api_errors.NotFoundError{Id: id}
	}
	return nil
}
func (ur *userRepository) GetAllUsers(page int, page_size int) (*[]models.User, error) {
	offset := page_size * (page - 1)

	query := "SELECT id, name FROM users ORDER BY id LIMIT $1 OFFSET $2"
	rows, err := ur.db.Query(query, page_size, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &users, nil
}
