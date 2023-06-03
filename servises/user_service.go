package servises

import models "game_leaderboards/m/v2/app/models"

type UserService interface {
	GetUsers() ([]models.User, error)
	GetUserByID(id int) (models.User, error)
	CreateUser(user models.User) error
	DeleteUser(id int) error
}
