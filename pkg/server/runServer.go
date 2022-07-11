package server

import (
	"os"
)

func RunServer() error {
	//if err := godotenv.Load("./.env"); err != nil {
	//	log.Fatal(err)
	//}
	server := newServer(os.Getenv("APP_HOST"), os.Getenv("APP_PORT"))
	return server.run()
}
