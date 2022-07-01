package checkStatus

import (
	"github.com/gin-gonic/gin"
	"github.com/rcrespodev/user_manager/api"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
)

func StatusGinHandlerFunc() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(200, &api.QueryResponse{
			Message: message.MessageData{},
			Data:    "Check-Status = Ok",
		})
	}
}
