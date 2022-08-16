package deleteUser

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rcrespodev/user_manager/api"
	"github.com/rcrespodev/user_manager/api/v1/handlers"
	delete "github.com/rcrespodev/user_manager/pkg/app/user/application/commands/delete"
	"github.com/rcrespodev/user_manager/pkg/kernel"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
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

		cmdUuid := uuid.New()
		log := domain.NewReturnLog(cmdUuid, kernel.Instance.MessageRepository(), "user")
		tokenUuid := ctx.GetString("token_uuid")
		if tokenUuid == "" || tokenUuid != clientArgs.UserUuid {
			log.LogError(domain.NewErrorCommand{
				NewMessageCommand: &message.NewMessageCommand{
					ObjectId:   clientArgs.UserUuid,
					MessageId:  0,
					MessagePkg: "Authorization",
				},
			})
			response = api.NewCommandResponse(log)
			ctx.JSON(int(log.HttpCode()), response)
			return
		}

		// if user is authorized, delete user command
		deleteUserCommand := delete.NewDeleteUserCommand(clientArgs.UserUuid)
		kernel.Instance.CommandBus().Exec(deleteUserCommand, log)

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
