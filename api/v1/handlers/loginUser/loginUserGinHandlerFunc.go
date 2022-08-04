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
)

func LoginUserGinHandlerFunc() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var clientArgs login.ClientArgs
		var response *api.CommandResponse
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

		response = api.NewCommandResponse(log)
		ctx.Set("jwt_key", cmdUuid.String())
		handlers.GinResponse(handlers.GinResponseCommand{
			Ctx:        ctx,
			Log:        log,
			StatusCode: int(log.HttpCode()),
			Data:       response,
		})
	}
}
