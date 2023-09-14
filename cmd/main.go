package main

import (
	jwtauth_app "github.com/Denialll/jwtauth-app"
	"github.com/Denialll/jwtauth-app/pkg/handler"
	"github.com/Denialll/jwtauth-app/pkg/repository"
	"github.com/Denialll/jwtauth-app/pkg/service"
	"github.com/spf13/viper"
	"log"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error init")
	}
	repos := repository.NewRepository()
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(jwtauth_app.Server)
	if err := srv.Run(viper.GetString("portc"), handlers.InitRoutes()); err != nil {
		log.Fatalf("error: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
