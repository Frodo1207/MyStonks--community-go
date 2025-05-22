package repository

import "MyStonks-go/internal/modules/users/models"

type UserRepositoryInterface interface {
	GetAll() ([]models.User, error)
	GetByID(id string) (*models.User, error)
	Create(user *models.User) error
}
