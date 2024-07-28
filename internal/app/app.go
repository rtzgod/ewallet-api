package app

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/rtzgod/EWallet/internal/domain/service"
	handlerHttp "github.com/rtzgod/EWallet/internal/handlers/http"
	"github.com/rtzgod/EWallet/internal/repository"
	"github.com/rtzgod/EWallet/internal/repository/postgres"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	if err := initConfig(); err != nil {
		logrus.Fatalf("Error initializing config file: %s", err)
	}
	if err := godotenv.Load(".env"); err != nil {
		logrus.Fatalf("Error loading .env file: %s", err)
	}

	db, err := postgres.NewPostgres(postgres.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		User:     viper.GetString("db.user"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   viper.GetString("db.name"),
		SSLMode:  viper.GetString("db.sslmode"),
	})

	if err != nil {
		logrus.Fatalf("Error initializing DB connection: %s", err)
	}

	repos := repository.NewRepository(db)

	services := service.NewService(repos)

	handlers := handlerHttp.NewHandler(services)

	server := new(Server)

	go func() {
		if err := server.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("Error occured while running server: %s", err)
		}
	}()

	logrus.Print("Server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Print("Shutting down server...")

	if err := server.Shutdown(context.Background()); err != nil {
		logrus.Errorf("Error occured while shutting down server: %s", err)
	}
	if err := db.Close(); err != nil {
		logrus.Errorf("Error occured while closing DB connection: %s", err)
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
