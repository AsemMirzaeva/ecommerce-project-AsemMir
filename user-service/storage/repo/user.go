package repo

import "user-service/models"

type UserStorage interface {
	CreateUser(user *models.User) (*models.User, error)
	GetUser(id string) (*models.User, error)
	UpdateUser(user *models.User) (*models.User, error)
	DeleteUser(id string) error
	ListUsers(limit, page int) ([]*models.User, error)
}
