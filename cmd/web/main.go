package main

import (
	"github.com/rcrespodev/user_manager/pkg/server"
	"log"
)

func main() {
	if err := server.RunServer(); err != nil {
		log.Fatal(err)
	}
}
