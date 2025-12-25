package middleware

import (
	"go-expense-management-system/internal/messages"
	"go-expense-management-system/internal/model"
	"go-expense-management-system/internal/usecase"
	"go-expense-management-system/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewAuth(userUserCase *usecase.UserUseCase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &model.VerifyUserRequest{Token: ctx.GetHeader("Authorization")}
		userUserCase.Log.Debugf("Authorization : %s", request.Token)

		auth, err := userUserCase.Verify(ctx.Request.Context(), request)
		if err != nil {
			res := utils.FailedResponse(messages.InvalidToken)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		userUserCase.Log.Debugf("User : %+v", auth.UserID)
		ctx.Set("auth", auth)
		ctx.Next()
	}
}

func GetUser(ctx *gin.Context) (*model.Auth, bool) {
	auth, ok := ctx.Get("auth")
	if !ok {
		return nil, false
	}

	typed, ok := auth.(*model.Auth)
	if !ok {
		return nil, false
	}

	return typed, true
}
