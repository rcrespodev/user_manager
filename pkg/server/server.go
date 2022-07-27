package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	apiEndpoints "github.com/rcrespodev/user_manager/api/v1/endpoints"
	"github.com/rcrespodev/user_manager/api/v1/handlers/jwtAuth"
	kernel "github.com/rcrespodev/user_manager/pkg/kernel"
	"log"
)

var Server *server

type server struct {
	httpAddress string
	engine      *gin.Engine
	kernel      *kernel.Kernel
}

func newServer(host, port string) *server {
	if Server != nil {
		return Server
	}
	Server = &server{
		httpAddress: fmt.Sprintf("%s:%s", host, port),
		engine:      gin.New(),
		kernel:      kernel.NewPrdKernel(nil, nil),
	}
	Server.engine.Use(jwtAuth.ValidateJwt()) //Jwt Auth
	return Server
}

func (s *server) run() error {
	routes := apiEndpoints.NewEndpoints()
	for _, route := range routes.Endpoints {
		s.engine.Handle(route.HttpMethod, route.RelativePath, route.Handler)
	}

	log.Printf("Server running on %v", s.httpAddress)
	return s.engine.Run(s.httpAddress)
}

func (s *server) Kernel() *kernel.Kernel {
	return s.kernel
}
