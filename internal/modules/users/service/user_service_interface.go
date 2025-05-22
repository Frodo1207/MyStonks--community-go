package service

import "MyStonks-go/internal/modules/users/models"

type UserServiceInterface interface {
	GetUsers() ([]models.User, error)
	GetUser(id string) (*models.User, error)
	CreateUser(user *models.CreateUserRequest) (*models.User, error)
}
