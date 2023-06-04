package repositories

import (
	"database/sql"

	models "game_leaderboards/m/v2/app/models"
)

type leaderboardRepository struct {
	db *sql.DB
}

// LeaderboardRepository представляет интерфейс для работы с пользователями в базе данных.
type LeaderboardRepository interface {
	GetGameLeaderboard(gameId int) (*models.GameLeaderboard, error)
	GetUserLeaderboard(userId int) (*[]models.UserLeaderboard, error)
	AddRecord(models.LeaderboardEntry) error
}

func NewLeaderboardRepository(db *sql.DB) LeaderboardRepository {
	return &leaderboardRepository{
		db: db,
	}
}

func (lr *leaderboardRepository) GetGameLeaderboard(gameId int) (*models.GameLeaderboard, error) {
	var gameName string
	gameQuery := "SELECT name FROM games WHERE id = $1"
	err := lr.db.QueryRow(gameQuery, gameId).Scan(&gameName)
	if err != nil {
		return nil, err
	}

	query := "SELECT users.name, score from leaderboard LEFT JOIN users ON users.id = user_id WHERE game_id = $1 ORDER BY score DESC;"
	rows, err := lr.db.Query(query, gameId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	leaderboards := []models.LeaderboardResponse{}
	for rows.Next() {
		var leaderboard models.LeaderboardResponse
		if err := rows.Scan(&leaderboard.UserName, &leaderboard.Score); err != nil {
			return nil, err
		}
		leaderboards = append(leaderboards, leaderboard)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	gameLeaderboard := models.GameLeaderboard{
		GameName:    gameName,
		Leaderboard: leaderboards,
	}

	return &gameLeaderboard, nil
}
func (lr *leaderboardRepository) AddRecord(entry models.LeaderboardEntry) error {
	query := "INSERT INTO leaderboard (gameId, userId, score) VALUES ($1, $2, $3)"
	_, err := lr.db.Exec(query, entry.GameID, entry.UserID, entry.Score)
	if err != nil {
		return err
	}

	return nil
}
func (lr *leaderboardRepository) GetUserLeaderboard(userId int) (*[]models.UserLeaderboard, error) {
	query := `SELECT g.name AS game_name, l.score, 
	(SELECT COUNT(*) + 1 
	 FROM leaderboard 
	 WHERE game_id = l.game_id AND score > l.score) AS position
FROM leaderboard l
JOIN games g ON l.game_id = g.id
WHERE l.user_id = $1;`

	rows, err := lr.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	leaderboards := []models.UserLeaderboard{}
	for rows.Next() {
		var leaderboard models.UserLeaderboard
		if err := rows.Scan(&leaderboard.GameName, &leaderboard.Score, &leaderboard.Position); err != nil {
			return nil, err
		}
		leaderboards = append(leaderboards, leaderboard)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &leaderboards, nil
}
