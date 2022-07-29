package jwtAuth

import (
	"github.com/gin-gonic/gin"
	"github.com/rcrespodev/user_manager/api/v1/endpoints"
	jwtDomain "github.com/rcrespodev/user_manager/pkg/app/auth/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel"
	"net/http"
)

func ValidateJwt() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if path == endpoints.EndpointRegisterUser || path == endpoints.EndpointCheckStatus {
			ctx.Next()
			return
		}

		tokenString := ctx.GetHeader("Authorization")
		if err := jwtDomain.IsValidJwt(tokenString, kernel.Instance.JwtConfig()); err != nil {
			ctx.JSON(401, nil)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Next()
	}
}

//func GenerateJwt(uuid uuid.UUID) gin.HandlerFunc {
//	return func(ctx *gin.Context) {
//		token, err := jwtDomain.SignJwt(uuid, kernel.Instance.JwtConfig())
//		if err != nil {
//			ctx.JSON(401, nil)
//			return
//		}
//		ctx.Header("Token", token)
//		ctx.Next()
//	}
//
//}
