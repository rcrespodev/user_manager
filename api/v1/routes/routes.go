package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rcrespodev/user_manager/api/v1/handlers/checkStatus"
	"github.com/rcrespodev/user_manager/api/v1/handlers/registerUser"
	"net/http"
)

type Routes struct {
	Routes []Route
}

type Route struct {
	HttpMethod   string
	RelativePath string
	Handler      gin.HandlerFunc
}

//func (r Route) RelativePath() string {
//	return r.RelativePath
//}
//
//func (r Route) Handler() gin.HandlerFunc {
//	return r.Handler
//}
//
//func (r Route) HttpMethod() string {
//	return r.HttpMethod
//}

func NewRoutes() Routes {
	routes := []Route{
		{
			HttpMethod:   http.MethodGet,
			RelativePath: "/check-status",
			Handler:      checkStatus.StatusGinHandlerFunc(),
		},
		{
			HttpMethod:   http.MethodPost,
			RelativePath: "/register_user",
			Handler:      registerUser.RegisterUserGinHandlerFunc(),
		},
	}
	return Routes{Routes: routes}
}
