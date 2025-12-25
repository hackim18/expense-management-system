package usecase

import (
	"context"
	"go-expense-management-system/internal/entity"
	"go-expense-management-system/internal/messages"
	"go-expense-management-system/internal/model"
	"go-expense-management-system/internal/model/converter"
	"go-expense-management-system/internal/repository"
	"go-expense-management-system/internal/utils"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUseCase struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	JWT            *utils.JWTHelper
	UserRepository *repository.UserRepository
}

func NewUserUseCase(db *gorm.DB, logger *logrus.Logger, jwt *utils.JWTHelper,
	userRepository *repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		DB:             db,
		Log:            logger,
		JWT:            jwt,
		UserRepository: userRepository,
	}
}

func (c *UserUseCase) Verify(ctx context.Context, request *model.VerifyUserRequest) (*model.Auth, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if request.Token == "" {
		return nil, model.ErrUnauthorized
	}

	tokenStr := strings.TrimPrefix(request.Token, "Bearer ")
	claims, err := c.JWT.DecodeAccessToken(tokenStr)
	if err != nil {
		c.Log.Warnf("Failed to decode access token : %+v", err)
		return nil, utils.Error(messages.InvalidToken, http.StatusUnauthorized, err)
	}

	if claims.Subject == "" {
		c.Log.Warnf("Invalid user_id in token claims")
		return nil, utils.Error(messages.InvalidToken, http.StatusUnauthorized, nil)
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		c.Log.Warnf("Invalid UUID format in token claims")
		return nil, utils.Error(messages.InvalidToken, http.StatusUnauthorized, err)
	}

	return &model.Auth{
		UserID: userID,
	}, nil
}

func (c *UserUseCase) Create(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	total, err := c.UserRepository.CountByCondition(tx, "email = ?", request.Email)
	if err != nil {
		c.Log.Warnf("Failed to check existing user : %+v", err)
		return nil, utils.Error(messages.ErrCheckUser, http.StatusInternalServerError, err)
	}

	if total > 0 {
		return nil, utils.Error(messages.ErrUserAlreadyExists, http.StatusConflict, nil)
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Log.Warnf("Failed to generate bcrypt hash : %+v", err)
		return nil, utils.Error(messages.ErrProcessPassword, http.StatusInternalServerError, err)
	}

	user := &entity.User{
		Name:         request.Name,
		Email:        request.Email,
		PasswordHash: string(password),
	}

	if err := c.UserRepository.Create(tx, user); err != nil {
		c.Log.Warnf("Failed to insert user : %+v", err)
		return nil, utils.Error(messages.ErrCreateUser, http.StatusInternalServerError, err)
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, utils.Error(messages.ErrCommitTransaction, http.StatusInternalServerError, err)
	}

	return converter.UserToResponse(user), nil
}

func (c *UserUseCase) Login(ctx context.Context, request *model.LoginUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	user := new(entity.User)
	if err := c.UserRepository.FindByCondition(tx, user, "email = ?", request.Email); err != nil {
		c.Log.Warnf("Failed to find user by email : %+v", err)
		return nil, utils.Error(messages.ErrInvalidEmailOrPassword, http.StatusUnauthorized, err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password)); err != nil {
		c.Log.Warnf("Invalid password : %+v", err)
		return nil, utils.Error(messages.ErrInvalidEmailOrPassword, http.StatusUnauthorized, err)
	}

	if c.JWT == nil {
		c.Log.Warn("JWT helper not configured")
		return nil, utils.Error(messages.ErrGenerateAccessToken, http.StatusInternalServerError, nil)
	}

	accessToken, err := c.JWT.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		c.Log.Warnf("Failed to generate access token : %+v", err)
		return nil, utils.Error(messages.ErrGenerateAccessToken, http.StatusInternalServerError, err)
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, model.ErrInternalServerError
	}

	return converter.UserToLoginResponse(user, accessToken), nil
}
