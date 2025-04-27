package main

import (
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

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

	err = godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		log.Fatal(err)
	}

	repos := repository.NewRepository(db)
	service := service.NewService(repos)
	handlers := handler.NewHandler(service)

	srv := &server.Server{}
	err = srv.Run(viper.GetString("8000"), handlers.InitRoute())
	if err != nil {
		log.Fatal(err)
	}
}
