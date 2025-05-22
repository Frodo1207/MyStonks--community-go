package repository

import (
	"MyStonks-go/internal/modules/users/models"
	"github.com/sirupsen/logrus"
)

type UserRepository struct {
	logger *logrus.Logger
	// 这里可以添加数据库连接等依赖
}

func NewUserRepository(logger *logrus.Logger) *UserRepository {
	return &UserRepository{
		logger: logger,
	}
}

func (r *UserRepository) GetAll() ([]models.User, error) {
	r.logger.Info("Fetching all users from repository")
	// 实际项目中这里会从数据库获取数据
	return []models.User{
		{ID: "1", Username: "user1", Email: "user1@example.com"},
		{ID: "2", Username: "user2", Email: "user2@example.com"},
	}, nil
}

func (r *UserRepository) GetByID(id string) (*models.User, error) {
	r.logger.WithField("user_id", id).Info("Fetching user by ID from repository")
	// 实际项目中这里会从数据库获取数据
	return &models.User{
		ID:       id,
		Username: "user" + id,
		Email:    "user" + id + "@example.com",
	}, nil
}

func (r *UserRepository) Create(user *models.User) error {
	r.logger.WithField("username", user.Username).Info("Creating user in repository")
	// 实际项目中这里会保存到数据库
	return nil
}
