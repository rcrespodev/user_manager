package endpoints

import (
	"github.com/gin-gonic/gin"
	"github.com/rcrespodev/user_manager/api/v1/handlers/checkStatus"
	"github.com/rcrespodev/user_manager/api/v1/handlers/getUser"
	"github.com/rcrespodev/user_manager/api/v1/handlers/logOutUser"
	"github.com/rcrespodev/user_manager/api/v1/handlers/loginUser"
	"github.com/rcrespodev/user_manager/api/v1/handlers/registerUser"
	"net/http"
)

const (
	EndpointCheckStatus  = "/v1/check-status"
	EndpointRegisterUser = "/v1/register_user"
	EndpointLogin        = "/v1/login"
	EndpointLogOut       = "/v1/logout"
	EndpointGetUser      = "/v1/user"
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
		EndpointGetUser: Endpoint{
			HttpMethod:     http.MethodGet,
			Handler:        getUser.GetUserGinHandlerFunc(),
			AuthValidation: true,
		},
	}
}
