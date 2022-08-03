package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rcrespodev/user_manager/api"
	"github.com/rcrespodev/user_manager/pkg/app/auth-jwt/application/commands/userLogged"
	jwtDomain "github.com/rcrespodev/user_manager/pkg/app/auth-jwt/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/command"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"net/http"
)

type GinResponseCommand struct {
	Ctx        *gin.Context
	Log        *domain.ReturnLog
	StatusCode int
	Data       interface{}
}

func GinResponse(cmd GinResponseCommand) {
	jwtKey := cmd.Ctx.GetString("jwt_key")
	if jwtKey == "" {
		cmd.Ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	jwtKeyUuid, err := uuid.Parse(jwtKey)
	if err != nil {
		cmd.Ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	userLoggedCommand := userLogged.NewCommand(jwtKeyUuid)
	c := command.NewCommand(command.UserLogged, jwtKeyUuid, userLoggedCommand)
	userLoggedLog := domain.NewReturnLog(jwtKeyUuid, kernel.Instance.MessageRepository(), "authorization")
	kernel.Instance.CommandBus().Exec(*c, userLoggedLog)
	if userLoggedLog.Error() != nil {
		response := api.NewCommandResponse(userLoggedLog)
		cmd.Ctx.JSON(int(userLoggedLog.HttpCode()), response)
		return
	}

	token := kernel.Instance.JwtRepository().FindByUuid(jwtDomain.FindByUuidQuery{Uuid: jwtKey})
	if token == nil {
		cmd.Ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if token.Token == "" {
		cmd.Ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	cmd.Ctx.Header("Token", token.Token)
	cmd.Ctx.JSON(cmd.StatusCode, cmd.Data)
}
