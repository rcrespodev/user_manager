package jwtAuth

import (
	"github.com/gin-gonic/gin"
	jwtDomain "github.com/rcrespodev/user_manager/pkg/app/auth/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel"
)

func ValidateJwt() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if path := ctx.Request.URL.Path; path == "/register_user" {
			ctx.Next()
		}

		tokenString := ctx.GetHeader("Authorization")
		if err := jwtDomain.ParseJwt(tokenString, kernel.Instance.JwtConfig()); err != nil {
			ctx.JSON(401, nil)
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
