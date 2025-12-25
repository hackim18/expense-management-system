package converter

import (
	"go-expense-management-system/internal/entity"
	"go-expense-management-system/internal/model"
)

func UserToResponse(user *entity.User) *model.UserResponse {
	id := user.ID
	return &model.UserResponse{
		ID:    &id,
		Name:  user.Name,
		Email: user.Email,
	}
}

func UserToLoginResponse(user *entity.User, accessToken string) *model.UserResponse {
	id := user.ID
	return &model.UserResponse{
		ID:          &id,
		Name:        user.Name,
		Email:       user.Email,
		AccessToken: accessToken,
	}
}
