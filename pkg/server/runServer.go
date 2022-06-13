package server

import "os"

func RunServer() error {
	server := newServer(os.Getenv("APP_HOST"), os.Getenv("APP_PORT"))
	return server.run()
}
