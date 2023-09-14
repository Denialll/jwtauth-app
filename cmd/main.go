package main

import (
	jwtauth_app "github.com/Denialll/jwtauth-app"
	"github.com/Denialll/jwtauth-app/pkg/handler"
	"github.com/Denialll/jwtauth-app/pkg/repository"
	"github.com/Denialll/jwtauth-app/pkg/service"
	"log"
)

func main() {
	repos := repository.NewRepository()
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	handlers = new(handler.Handler)
	srv := new(jwtauth_app.Server)
	if err := srv.Run("8080", handlers.InitRoutes()); err != nil {
		log.Fatalf("error: %s", err.Error())
	}
}
