package endpoints

import (
	"fmt"
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
	EndpointCheckStatus = "/v1/check-status"
	EndpointUser        = "/v1/user/"
	EndpointUserLogin   = "/v1/user/login/"
)

type Endpoints map[string]Endpoint

type Endpoint struct {
	HttpMethod     string
	RelPath        string
	Handler        gin.HandlerFunc
	AuthValidation bool
}

func NewEndpoints() Endpoints {
	return Endpoints{
		fmt.Sprintf("%s-%s", http.MethodGet, EndpointCheckStatus): Endpoint{
			HttpMethod:     http.MethodGet,
			RelPath:        EndpointCheckStatus,
			Handler:        checkStatus.StatusGinHandlerFunc(),
			AuthValidation: false,
		},
		// user
		fmt.Sprintf("%s-%s", http.MethodPost, EndpointUser): Endpoint{
			HttpMethod:     http.MethodPost,
			RelPath:        EndpointUser,
			Handler:        registerUser.RegisterUserGinHandlerFunc(),
			AuthValidation: false,
		},
		fmt.Sprintf("%s-%s", http.MethodDelete, EndpointUser): Endpoint{
			HttpMethod:     http.MethodDelete,
			RelPath:        EndpointUser,
			Handler:        deleteUser.DeleteUserGinHandlerFunc(),
			AuthValidation: true,
		},
		fmt.Sprintf("%s-%s", http.MethodGet, EndpointUser): Endpoint{
			HttpMethod:     http.MethodGet,
			RelPath:        EndpointUser,
			Handler:        getUser.GetUserGinHandlerFunc(),
			AuthValidation: true,
		},
		// login logout user
		fmt.Sprintf("%s-%s", http.MethodPost, EndpointUserLogin): Endpoint{
			HttpMethod:     http.MethodPost,
			RelPath:        EndpointUserLogin,
			Handler:        loginUser.LoginUserGinHandlerFunc(),
			AuthValidation: false,
		},
		fmt.Sprintf("%s-%s", http.MethodDelete, EndpointUserLogin): Endpoint{
			HttpMethod:     http.MethodDelete,
			RelPath:        EndpointUserLogin,
			Handler:        logOutUser.LogOutUserGinHandlerFunc(),
			AuthValidation: true,
		},
	}
}
