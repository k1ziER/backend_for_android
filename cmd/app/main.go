package main

import (
	"android/internal/delivery/http/handler"
	"android/internal/repository"
	"android/internal/server"
	"android/internal/service"
	"log"
)

func main() {
	repos := repository.NewRepository()
	service := service.NewService(repos)
	handlers := handler.NewHandler(service)

	srv := &server.Server{}
	err := srv.Run("8080", handlers.InitRoute())
	if err != nil {
		log.Fatal(err)
	}
}
