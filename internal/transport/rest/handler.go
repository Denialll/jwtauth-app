package handler

import (
	_ "github.com/Denialll/jwtauth-app/docs"
	"github.com/Denialll/jwtauth-app/internal/service"
	"github.com/Denialll/jwtauth-app/pkg"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services     *service.Service
	tokenManager pkg.TokenManager
}

func NewHandler(services *service.Service, tokenManager pkg.TokenManager) *Handler {
	return &Handler{
		services:     services,
		tokenManager: tokenManager,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.POST("/refresh", h.refresh)
	}

	api := router.Group("/checkjwt", h.userIdentity)
	{
		api.GET("/", h.checkJWT)
	}

	return router
}
