package logOutUser

import (
	"github.com/gin-gonic/gin"
)

func LogOutUserGinHandlerFunc() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//tokenString := ctx.GetHeader("Authorization")
		//kernel.Instance.CommandBus().Exec()
	}
}
