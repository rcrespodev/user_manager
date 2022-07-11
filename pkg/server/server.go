package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	apiRoutes "github.com/rcrespodev/user_manager/api/v1/routes"
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
		kernel:      kernel.NewPrdKernel(),
	}
	return Server
}

func (s *server) run() error {
	routes := apiRoutes.NewRoutes()
	for _, route := range routes.Routes {
		s.engine.Handle(route.HttpMethod, route.RelativePath, route.Handler)
	}

	log.Printf("Server running on %v", s.httpAddress)
	return s.engine.Run(s.httpAddress)
}

func (s *server) Kernel() *kernel.Kernel {
	return s.kernel
}
