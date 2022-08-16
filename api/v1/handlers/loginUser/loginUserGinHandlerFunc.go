package loginUser

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rcrespodev/user_manager/api"
	"github.com/rcrespodev/user_manager/api/v1/handlers"
	"github.com/rcrespodev/user_manager/pkg/app/user/application/commands/login"
	"github.com/rcrespodev/user_manager/pkg/app/user/application/querys/userFinder"
	userDomain "github.com/rcrespodev/user_manager/pkg/app/user/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/valueObjects"
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
		cmdBus := kernel.Instance.CommandBus()
		cmdBus.Exec(loginCommand, log)

		if log.Status() == valueObjects.Success {
			queryArgs := []userDomain.FindUserQuery{
				{
					Log: log,
					Where: []userDomain.WhereArgs{
						{
							Field: "alias",
							Value: loginCommand.AliasOrEmail(),
						},
					},
				},
				{
					Log: log,
					Where: []userDomain.WhereArgs{
						{
							Field: "email",
							Value: loginCommand.AliasOrEmail(),
						},
					},
				},
			}
			FindUserQuery := userFinder.NewQuery(queryArgs)
			userSchema := kernel.Instance.QueryBus().Exec(FindUserQuery, log)
			if userSchema == nil {
				ctx.AbortWithStatus(500)
				return
			}
			user, ok := userSchema.(*userDomain.UserSchema)
			if !ok {
				ctx.AbortWithStatus(500)
				return
			}
			ctx.Set("jwt_key", user.Uuid)
		}

		response = api.NewCommandResponse(log)
		handlers.GinResponse(handlers.GinResponseCommand{
			Ctx:        ctx,
			Log:        log,
			StatusCode: int(log.HttpCode()),
			Data:       response,
		})
	}
}
