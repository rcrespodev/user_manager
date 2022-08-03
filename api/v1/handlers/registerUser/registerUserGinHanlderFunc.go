package registerUser

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rcrespodev/user_manager/api"
	"github.com/rcrespodev/user_manager/api/v1/handlers"
	"github.com/rcrespodev/user_manager/pkg/app/user/application/commands/register"
	"github.com/rcrespodev/user_manager/pkg/kernel"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/command"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"time"
)

func RegisterUserGinHandlerFunc() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var clientArgs register.ClientArgs
		var response *api.CommandResponse
		if err := ctx.BindJSON(&clientArgs); err != nil {
			response.Message = message.MessageData{
				ObjectId:        "",
				MessageId:       1,
				MessagePkg:      "http handler",
				Variables:       message.Variables{},
				Text:            "body request has´t correct type",
				Time:            time.Now(),
				ClientErrorType: message.ClientErrorBadRequest,
			}
			ctx.JSON(400, response)
			return
		}
		cmdUuid, err := uuid.Parse(clientArgs.Uuid)
		if err != nil {
			response.Message = message.MessageData{
				ObjectId:        "",
				MessageId:       1,
				MessagePkg:      "http handler",
				Variables:       message.Variables{},
				Text:            "body uuid has´t correct type",
				Time:            time.Now(),
				ClientErrorType: message.ClientErrorBadRequest,
			}
			ctx.JSON(400, response)
			return
		}

		log := domain.NewReturnLog(cmdUuid, kernel.Instance.MessageRepository(), "user")

		registerUserCommand := register.NewRegisterUserCommand(clientArgs)
		cmd := command.NewCommand(command.RegisterUser, cmdUuid, registerUserCommand)

		cmdBus := kernel.Instance.CommandBus()
		cmdBus.Exec(*cmd, log)

		//switch log.Status() {
		//case valueObjects.Error:
		//	if log.Error().InternalError() != nil {
		//		response.Message = message.MessageData{}
		//	} else {
		//		response.Message = *log.Error().Message()
		//	}
		//case valueObjects.Success:
		//	response.Message = *log.Success().MessageData()
		//}
		response = api.NewCommandResponse(log)
		//ctx.JSON(int(log.HttpCode()), response)
		ctx.Set("jwt_key", cmdUuid.String())
		handlers.GinResponse(handlers.GinResponseCommand{
			Ctx:        ctx,
			Log:        log,
			StatusCode: int(log.HttpCode()),
			Data:       response,
		})

		//token, err := jwtDomain.SignJwt(cmdUuid, time.Now(), kernel.Instance.Jwt())
		//if err != nil {
		//	ctx.JSON(401, nil)
		//	return
		//} else {
		//	if log.HttpCode() == 200 {
		//		ctx.Header("Token", token)
		//	}
		//	ctx.JSON(int(log.HttpCode()), response)
		//}
	}
}
