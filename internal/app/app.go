package app

import (
	"context"
	"github.com/Denialll/jwtauth-app/internal/repository"
	"github.com/Denialll/jwtauth-app/internal/service"
	handler "github.com/Denialll/jwtauth-app/internal/transport/rest"
	"github.com/Denialll/jwtauth-app/pkg"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
)

// @title JWTauth API
// @version 1.0
// @description REST API for JWTauth(Refresh + Access tokens)

// @host localhost:8085
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func Run() {
	gin.SetMode(gin.ReleaseMode)

	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err != nil {
		logrus.Fatalf("error init config: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error env: %s", err.Error())
	}

	tokenManager, err := pkg.NewManager(os.Getenv("JWT_KEY"), viper.GetDuration("jwt.accessTokenTTL"), viper.GetDuration("jwt.refreshTokenTTL"))
	if err != nil {
		logrus.Error(err)
		return
	}

	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		logrus.Fatalf("error MongoDB: %s", err.Error())
	}

	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Print(err)
		}
	}()

	db := mongoClient.Database(os.Getenv("MONGO_DB_NAME"))

	repos := repository.NewRepository(db)
	services := service.NewService(service.Deps{
		Repos:        repos,
		TokenManager: tokenManager,
		AccessTTL:    viper.GetDuration("jwt.accessTokenTTL"),
		RefreshTTL:   viper.GetDuration("jwt.refreshTokenTTL"),
	})
	handlers := handler.NewHandler(services, tokenManager)

	srv := &http.Server{
		Addr:    viper.GetString("port"),
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
