package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"android/internal/config"
	"android/internal/delivery/http/handler"
	"android/internal/repository"
	"android/internal/server"
	"android/internal/service"

	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	err := config.InitConfig()
	if err != nil {
		logrus.Fatal(err)
	}

	err = godotenv.Load()
	if err != nil {
		logrus.Fatal(err)
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
		logrus.Fatal(err)
	}

	repos := repository.NewRepository(db)
	service := service.NewService(repos)
	handlers := handler.NewHandler(service)

	fmt.Println("Start server at :8080")
	srv := &server.Server{}

	err = srv.Run(viper.GetString("port"), handlers.InitRoute())
	if err != nil {
		fmt.Println("Start server at :8080")
		logrus.Fatal(err)
	}
	fmt.Println("Start server at :8080")
}
