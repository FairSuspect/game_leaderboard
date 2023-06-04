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

type gameRepository struct {
	db *sql.DB
}

// GameRepository представляет интерфейс для работы с пользователями в базе данных.
type GameRepository interface {
	GetAllGames() (*[]models.Game, error)
	GetGameByID(id int) (*models.Game, error)
	CreateGame(game models.Game) (*models.Game, error)
	EditGame(game models.Game) (*models.Game, error)
	DeleteGame(id int) error
}

func NewGameRepository(db *sql.DB) GameRepository {
	return &gameRepository{
		db: db,
	}
}

func (ur *gameRepository) CreateGame(game models.Game) (*models.Game, error) {
	rows := ur.db.QueryRow("INSERT INTO games (name) VALUES ($1) RETURNING ID", game.Name)

	var id int
	err := rows.Scan(&id)
	if err != nil {

		// Проверяем, является ли ошибка ошибкой нарушения уникального ограничения
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Code == "23505" { // код "23505" соответствует ошибке нарушения уникального ограничения
			// Формируем свое сообщение об ошибке
			errMsg := fmt.Sprintf("Gamename '%s' already exists", *game.Name)
			return nil, errors.New(errMsg)
		}

		// Возвращаем другие ошибки без изменений
		return nil, err

	}
	game.ID = int(id)
	return &game, nil
}
func (ur *gameRepository) EditGame(game models.Game) (*models.Game, error) {
	_, err := ur.db.Exec("UPDATE games SET name = $1 WHERE id = $2", game.Name, game.ID)
	if err != nil {

		// Возвращаем другие ошибки без изменений
		return nil, err

	}
	return &game, err
}

func (ur *gameRepository) GetGameByID(id int) (*models.Game, error) {
	var game models.Game
	err := ur.db.QueryRow("SELECT id, name FROM games WHERE id = $1", id).Scan(&game.ID, &game.Name)
	if err != nil {
		log.Println("Failed to get game by ID:", err)
		return nil, api_errors.NotFoundError{Id: id}
	}
	return &game, nil
}
func (ur *gameRepository) DeleteGame(id int) error {
	_, err := ur.db.Exec("DELETE FROM games WHERE id = $1", id)
	if err != nil {
		log.Println("Failed to get game by ID:", err)
		return api_errors.NotFoundError{Id: id}
	}
	return nil
}
func (ur *gameRepository) GetAllGames() (*[]models.Game, error) {
	query := "SELECT id, name FROM games ORDER BY id"
	rows, err := ur.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	games := []models.Game{}
	for rows.Next() {
		var game models.Game
		if err := rows.Scan(&game.ID, &game.Name); err != nil {
			return nil, err
		}
		games = append(games, game)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &games, nil
}
