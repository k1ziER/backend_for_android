package main

import (
	"android/internal/config"
	"android/internal/delivery/http/handler"
	"android/internal/repository"
	"android/internal/server"
	"android/internal/service"
	"log"

	"github.com/spf13/viper"
)

func main() {
	err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	repos := repository.NewRepository()
	service := service.NewService(repos)
	handlers := handler.NewHandler(service)

	srv := &server.Server{}
	err = srv.Run(viper.GetString("8000"), handlers.InitRoute())
	if err != nil {
		log.Fatal(err)
	}
}
