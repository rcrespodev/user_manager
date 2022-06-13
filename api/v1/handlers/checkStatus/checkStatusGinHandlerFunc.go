package checkStatus

import "github.com/gin-gonic/gin"

func StatusGinHandlerFunc() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(200, "Check-Status = Ok")
	}
}
