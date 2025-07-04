package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"android/internal/config"
	"android/internal/delivery/http/handler"
	"android/internal/redis"
	"android/internal/repository"
	"android/internal/server"
	"android/internal/service"

	"github.com/spf13/viper"
)

// @title AndroidBackend
// @version 1.0
// @description API Server for Android client

// @host localhost:8090
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
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
	adress := os.Getenv("redis")

	cacheTTL := time.Hour

	redisClient := redis.NewRedisClient(adress, "", 0, cacheTTL)
	defer redisClient.Close()

	// Проверка соединения
	if err := redisClient.Ping(context.Background()); err != nil {
		logrus.Fatalf("Failed to connect to Redis: %v", err)
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

	blacklist := service.NewTokenBlacklist()
	repos := repository.NewRepository(db, redisClient)
	service := service.NewService(repos, blacklist)
	handlers := handler.NewHandler(service, blacklist)

	srv := &server.Server{}
	go func() {
		err = srv.Run(viper.GetString("port"), handlers.InitRoute())
		if err != nil {
			logrus.Fatal(err)
		}
	}()

	logrus.Println("Server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT)
	<-quit

	logrus.Println("Server shutting down")
	err = srv.Shutdown(context.Background())
	if err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}
	err = db.Close()
	if err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}

}
