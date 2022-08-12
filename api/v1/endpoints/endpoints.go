package endpoints

import (
	"github.com/gin-gonic/gin"
	"github.com/rcrespodev/user_manager/api/v1/handlers/checkStatus"
	"github.com/rcrespodev/user_manager/api/v1/handlers/deleteUser"
	"github.com/rcrespodev/user_manager/api/v1/handlers/getUser"
	"github.com/rcrespodev/user_manager/api/v1/handlers/logOutUser"
	"github.com/rcrespodev/user_manager/api/v1/handlers/loginUser"
	"github.com/rcrespodev/user_manager/api/v1/handlers/registerUser"
	"net/http"
)

const (
	EndpointCheckStatus  = "/v1/check-status"
	EndpointRegisterUser = "/v1/user/register" // /v1/user/ post
	EndpointLogin        = "/v1/user/login"    // /v1/user/login/ post
	EndpointLogOut       = "/v1/user/logout"   // /v1/user/login/ delete
	EndpointGetUser      = "/v1/user"          // /v1/user/ get
	EndpointDeleteUser   = "/v1/user/delete"   // /v1/user/ delete
)

type Endpoints map[string]Endpoint

type Endpoint struct {
	HttpMethod     string
	Handler        gin.HandlerFunc
	AuthValidation bool
}

func NewEndpoints() Endpoints {
	return Endpoints{
		EndpointCheckStatus: Endpoint{
			HttpMethod:     http.MethodGet,
			Handler:        checkStatus.StatusGinHandlerFunc(),
			AuthValidation: false,
		},
		EndpointRegisterUser: Endpoint{
			HttpMethod:     http.MethodPost,
			Handler:        registerUser.RegisterUserGinHandlerFunc(),
			AuthValidation: false,
		},
		EndpointLogin: Endpoint{
			HttpMethod:     http.MethodPost,
			Handler:        loginUser.LoginUserGinHandlerFunc(),
			AuthValidation: false,
		},
		EndpointLogOut: Endpoint{
			HttpMethod:     http.MethodPost,
			Handler:        logOutUser.LogOutUserGinHandlerFunc(),
			AuthValidation: true,
		},
		EndpointDeleteUser: Endpoint{
			HttpMethod:     http.MethodPost,
			Handler:        deleteUser.DeleteUserGinHandlerFunc(),
			AuthValidation: true,
		},
		EndpointGetUser: Endpoint{
			HttpMethod:     http.MethodGet,
			Handler:        getUser.GetUserGinHandlerFunc(),
			AuthValidation: true,
		},
	}
}
