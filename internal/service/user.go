package service

import (
	"MyStonks-go/internal/models"
)

func CreateUserIfNotExists(solAddress string) error {
	return models.CreateUserIfNotExists(solAddress)
}

func GetUserInfo(solAddress string) (*models.User, error) {
	user, err := models.GetUserBySolAddress(solAddress)
	if err != nil {
		return nil, err
	}
	return user, nil
}
