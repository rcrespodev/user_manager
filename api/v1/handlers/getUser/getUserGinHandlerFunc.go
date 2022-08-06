package getUser

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rcrespodev/user_manager/api"
	"github.com/rcrespodev/user_manager/api/v1/handlers"
	"github.com/rcrespodev/user_manager/pkg/app/user/application/querys/userFinder"
	"github.com/rcrespodev/user_manager/pkg/app/user/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel"
	query "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/query"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"net/http"
)

func GetUserGinHandlerFunc() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cmdUuid := uuid.New()
		tokenUuid := ctx.GetString("token_uuid")
		retLog := returnLog.NewReturnLog(cmdUuid, kernel.Instance.MessageRepository(), "user")
		retLog.SetObjectId(tokenUuid)

		var queryArgs []domain.WhereArgs
		for _, allowedArg := range getAllowedArgs() {
			queryValue := ctx.Query(allowedArg)
			if queryValue != "" {
				queryArgs = append(queryArgs, domain.WhereArgs{
					Field: allowedArg,
					Value: queryValue,
				})
			}
		}

		findUserQuery := userFinder.NewQuery([]domain.FindUserQuery{
			{
				Log:   retLog,
				Where: queryArgs,
			},
		})
		q := query.NewQuery(query.FindUser, findUserQuery)
		data := kernel.Instance.QueryBus().Exec(q, retLog)

		userSchema, ok := data.(*domain.UserSchema)
		if !ok {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if userSchema != nil {
			userSchema.HashedPassword = []byte{}
			data = userSchema
			compareTokenUuidAndUserUuid(tokenUuid, userSchema.Uuid, retLog)
			ctx.Set("jwt_key", userSchema.Uuid)
		}

		response := api.NewQueryResponse(retLog, data)
		handlers.GinResponse(handlers.GinResponseCommand{
			Ctx:        ctx,
			Log:        retLog,
			StatusCode: int(retLog.HttpCode()),
			Data:       response,
		})

	}
}

func getAllowedArgs() []string {
	return []string{
		"uuid", "alias", "email", "name", "second_name",
	}
}

func compareTokenUuidAndUserUuid(tokenUuid, userUuid string, log *returnLog.ReturnLog) {
	if tokenUuid != userUuid {
		log.SetObjectId(tokenUuid)
		log.LogError(returnLog.NewErrorCommand{
			NewMessageCommand: &message.NewMessageCommand{
				MessageId:  0,
				MessagePkg: "Authorization",
			},
			Overwritten: true,
		})
	}
}
