package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	apiRoutes "github.com/rcrespodev/user_manager/api/v1/routes"
	"log"
)

type server struct {
	httpAddress string
	engine      *gin.Engine
}

func newServer(host, port string) *server {
	return &server{
		httpAddress: fmt.Sprintf("%s:%s", host, port),
		engine:      gin.New(),
	}
}

func (s *server) run() error {
	routes := apiRoutes.NewRoutes()
	for _, route := range routes.Routes {
		s.engine.Handle(route.HttpMethod, route.RelativePath, route.Handler)
	}

	log.Printf("Server running on %v", s.httpAddress)
	return s.engine.Run(s.httpAddress)
}
