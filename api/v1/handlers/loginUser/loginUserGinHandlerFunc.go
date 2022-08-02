package loginUser

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rcrespodev/user_manager/api"
	"github.com/rcrespodev/user_manager/api/v1/handlers"
	"github.com/rcrespodev/user_manager/pkg/app/user/application/commands/login"
	"github.com/rcrespodev/user_manager/pkg/kernel"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/command"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/valueObjects"
)

func LoginUserGinHandlerFunc() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var clientArgs login.ClientArgs
		var response api.CommandResponse
		if err := ctx.BindJSON(&clientArgs); err != nil {
			response.Message = handlers.BodyRequestBadType()
			ctx.JSON(400, response)
			return
		}
		cmdUuid := uuid.New()
		log := domain.NewReturnLog(cmdUuid, kernel.Instance.MessageRepository(), "user")
		loginCommand := login.NewLoginUserCommand(clientArgs)
		cmd := command.NewCommand(command.LoginUser, cmdUuid, loginCommand)
		cmdBus := kernel.Instance.CommandBus()
		cmdBus.Exec(*cmd, log)

		switch log.Status() {
		case valueObjects.Error:
			if log.Error().InternalError() != nil {
				response.Message = message.MessageData{}
			} else {
				response.Message = *log.Error().Message()
			}
		case valueObjects.Success:
			response.Message = *log.Success().MessageData()
		}
		ctx.JSON(int(log.HttpCode()), response)
		ctx.Set("jwt_key", cmdUuid)
		handlers.GinResponse(ctx)
	}
}
