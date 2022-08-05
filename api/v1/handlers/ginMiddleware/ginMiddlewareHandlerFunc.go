package ginMiddleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rcrespodev/user_manager/api/v1/endpoints"
	"github.com/rcrespodev/user_manager/pkg/kernel"
	"net/http"
)

func MiddlewareHandlerFunc() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestPath := endpoints.ParseEndpoint(ctx.Request.URL.Path)

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
		claims, err := kernel.Instance.Jwt().ValidateToken(tokenString)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Set("uuid", claims["key"])
		ctx.Next()
	}
}
