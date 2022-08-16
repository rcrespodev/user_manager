package logOutUser

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rcrespodev/user_manager/api"
	"github.com/rcrespodev/user_manager/pkg/app/auth-jwt/application/commands/userLoggedOut"
	"github.com/rcrespodev/user_manager/pkg/kernel"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"net/http"
)

func LogOutUserGinHandlerFunc() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		claims, err := kernel.Instance.Jwt().ValidateToken(tokenString)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenUuidStr, ok := claims["key"].(string)
		if !ok {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		tokenUuid, err := uuid.Parse(tokenUuidStr)
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		userLoggedOutCmd := userLoggedOut.NewCommand(tokenUuidStr)
		retLog := domain.NewReturnLog(tokenUuid, kernel.Instance.MessageRepository(), "authorization")
		kernel.Instance.CommandBus().Exec(userLoggedOutCmd, retLog)
		response := api.NewCommandResponse(retLog)
		ctx.JSON(int(retLog.HttpCode()), response)
	}
}
