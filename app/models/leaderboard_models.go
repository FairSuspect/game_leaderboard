package models

type LeaderboardEntry struct {
	GameID int `json:"gameId"`
	UserID int `json:"userId"`
	Score  int `json:"score"`
}

type LeaderboardResponse struct {
	UserName string `json:"userName"`
	Score    int    `json:"score"`
}

type UserLeaderboard struct {
	GameName string `json:"gameName"`
	Score    int    `json:"score"`
	Position int    `json:"position"`
}

type GameLeaderboard struct {
	GameName    string                `json:"gameName"`
	Leaderboard []LeaderboardResponse `json:"leaderboard"`
}
