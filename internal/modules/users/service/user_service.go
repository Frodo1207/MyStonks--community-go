package service

import (
	"MyStonks-go/internal/modules/users/models"
	"MyStonks-go/internal/modules/users/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type UserService struct {
	logger     *logrus.Logger
	repository repository.UserRepositoryInterface
}

func NewUserService(logger *logrus.Logger, repo repository.UserRepositoryInterface) *UserService {
	return &UserService{
		logger:     logger,
		repository: repo,
	}
}

func (s *UserService) GetUsers() ([]models.User, error) {
	s.logger.Info("Getting all users from service")
	return s.repository.GetAll()
}

func (s *UserService) GetUser(id string) (*models.User, error) {
	s.logger.WithField("user_id", id).Info("Getting user by ID from service")
	return s.repository.GetByID(id)
}

func (s *UserService) CreateUser(req *models.CreateUserRequest) (*models.User, error) {
	s.logger.WithField("username", req.Username).Info("Creating new user in service")

	user := &models.User{
		ID:       uuid.New().String(),
		Username: req.Username,
		Email:    req.Email,
	}

	if err := s.repository.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}
