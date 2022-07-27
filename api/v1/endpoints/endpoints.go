package endpoints

import (
	"github.com/gin-gonic/gin"
	"github.com/rcrespodev/user_manager/api/v1/handlers/checkStatus"
	"github.com/rcrespodev/user_manager/api/v1/handlers/loginUser"
	"github.com/rcrespodev/user_manager/api/v1/handlers/registerUser"
	"net/http"
)

const (
	EndpointCheckStatus  = "/v1/check-status"
	EndpointRegisterUser = "/v1/register_user"
	EndpointLogin        = "/v1/login"
)

type Endpoints struct {
	Endpoints []Endpoint
}

type Endpoint struct {
	HttpMethod   string
	RelativePath string
	Handler      gin.HandlerFunc
}

func NewEndpoints() Endpoints {
	endpoints := []Endpoint{
		{
			HttpMethod:   http.MethodGet,
			RelativePath: EndpointCheckStatus,
			Handler:      checkStatus.StatusGinHandlerFunc(),
		},
		{
			HttpMethod:   http.MethodPost,
			RelativePath: EndpointRegisterUser,
			Handler:      registerUser.RegisterUserGinHandlerFunc(),
		},
		{
			HttpMethod:   http.MethodPost,
			RelativePath: EndpointLogin,
			Handler:      loginUser.LoginUserGinHandlerFunc(),
		},
	}
	return Endpoints{Endpoints: endpoints}
}
