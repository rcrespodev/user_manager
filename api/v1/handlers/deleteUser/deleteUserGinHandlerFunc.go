package deleteUser

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rcrespodev/user_manager/api"
	"github.com/rcrespodev/user_manager/api/v1/handlers"
	delete "github.com/rcrespodev/user_manager/pkg/app/user/application/commands/delete"
	"github.com/rcrespodev/user_manager/pkg/kernel"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/command"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"net/http"
)

func DeleteUserGinHandlerFunc() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// First check if token uuid is the same of user uuid
		var clientArgs delete.ClientArgs
		var response *api.CommandResponse
		if err := ctx.BindJSON(&clientArgs); err != nil {
			response.Message = handlers.BodyRequestBadType()
			ctx.JSON(400, response)
			return
		}

		tokenUuid := ctx.GetString("token_uuid")
		if tokenUuid == "" || tokenUuid != clientArgs.UserUuid {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// if user is authorized, delete user command
		cmdUuid := uuid.New()

		log := domain.NewReturnLog(cmdUuid, kernel.Instance.MessageRepository(), "user")
		deleteUserCommand := delete.NewDeleteUserCommand(clientArgs.UserUuid)
		cmd := command.NewCommand(command.DeleteUser, cmdUuid, deleteUserCommand)
		kernel.Instance.CommandBus().Exec(*cmd, log)

		ctx.Set("jwt_key", deleteUserCommand.UserUuid())
		response = api.NewCommandResponse(log)
		handlers.GinResponse(handlers.GinResponseCommand{
			Ctx:        ctx,
			Log:        log,
			StatusCode: int(log.HttpCode()),
			Data:       response,
		})

	}
}
