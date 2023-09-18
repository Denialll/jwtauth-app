package main

import (
	repository2 "github.com/Denialll/jwtauth-app/internal/repository"
	"github.com/Denialll/jwtauth-app/internal/services"
	"github.com/Denialll/jwtauth-app/internal/transport/rest"
	"github.com/Denialll/jwtauth-app/pkg"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"os"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err != nil {
		logrus.Fatalf("error init config: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error env: %s", err.Error())
	}

	db, err := repository2.NewPostgresDB(repository2.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})

	if err != nil {
		logrus.Fatalf("Failed to init db: %s", err.Error())
	}

	tokenManager, err := pkg.NewManager("aaaaa", viper.GetDuration("jwt.accessTokenTTL"), viper.GetDuration("jwt.refreshTokenTTL"))
	if err != nil {
		logrus.Error(err)
		return
	}

	repos := repository2.NewRepository(db)
	services := services.NewService(services.Deps{
		Repos:        repos,
		TokenManager: tokenManager,
	})
	handlers := handler.NewHandler(services, tokenManager)

	srv := &http.Server{
		Addr:    ":8000",
		Handler: handlers.InitRoutes(),
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logrus.Fatalf("error: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("internal/config")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
