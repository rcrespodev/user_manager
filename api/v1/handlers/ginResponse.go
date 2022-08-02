package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/rcrespodev/user_manager/pkg/kernel"
)

func GinResponse(ctx *gin.Context) {
	jwtKey, ok := ctx.Get("jwt_key")
	if !ok {
		ctx.AbortWithStatus(500)
	}
	jwtKeyString, ok := jwtKey.(string)
	if !ok {
		ctx.AbortWithStatus(500)
	}
	token, err := kernel.Instance.Jwt().CreateNewToken(jwtKeyString)
	if err != nil {
		ctx.AbortWithStatus(500)
	}
	ctx.Header("Token", token)
}
