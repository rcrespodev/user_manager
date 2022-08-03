package main

import (
	"github.com/joho/godotenv"
	"github.com/rcrespodev/user_manager/pkg/server"
	"log"
)

func main() {
	err := godotenv.Load("/app/.env")
	if err != nil {
		log.Fatal(err)
	}
	if err := server.RunServer(); err != nil {
		log.Fatal(err)
	}
}
