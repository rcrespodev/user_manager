package ginMiddleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rcrespodev/user_manager/api/v1/endpoints"
	jwtDomain "github.com/rcrespodev/user_manager/pkg/app/auth/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel"
	"net/http"
)

func MiddlewareHandlerFunc() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestPath := ctx.Request.URL.Path

		// check path existence and method allowed
		appEndpoints := endpoints.NewEndpoints()
		endPoint, ok := appEndpoints[requestPath]
		if !ok {
			ctx.JSON(http.StatusNotFound, nil)
			return
		}
		if endPoint.HttpMethod != ctx.Request.Method {
			ctx.JSON(http.StatusMethodNotAllowed, nil)
			return
		}

		// check if endpoint not should execute JWT validation
		if !endPoint.AuthValidation {
			ctx.Next()
			return
		}

		// JWT validation
		tokenString := ctx.GetHeader("Authorization")
		if err := jwtDomain.ParseJwt(tokenString, kernel.Instance.JwtConfig()); err != nil {
			ctx.JSON(401, nil)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Next()
	}
}
