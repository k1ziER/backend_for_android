package main

import (
	"android/internal/delivery/http/handler"
	"android/internal/server"
	"log"
)

func main() {
	handlers := &handler.Handler{}
	srv := &server.Server{}
	err := srv.Run("8080", handlers.InitRoute())
	if err != nil {
		log.Fatal(err)
	}
}
