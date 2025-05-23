package service

import (
	"MyStonks-go/internal/models"
	"MyStonks-go/internal/service/schema"
)

func CreateUser(user *schema.UserCreateReq) (*schema.UserResp, error) {
	userModel := &models.User{
		Name:  user.Name,
		Email: user.Email,
	}
	err := models.CreateUser(userModel)
	if err != nil {
		return nil, err
	}
	return &schema.UserResp{
		ID:    userModel.ID,
		Name:  userModel.Name,
		Email: userModel.Email,
	}, nil
}

func GetUserByID(id string) (*schema.UserResp, error) {
	user, err := models.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	return &schema.UserResp{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
