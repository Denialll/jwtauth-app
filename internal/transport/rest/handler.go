package handler

import (
	"github.com/Denialll/jwtauth-app/internal/services"
	"github.com/Denialll/jwtauth-app/pkg"
	"github.com/gin-gonic/gin"
)

//type Handler struct {
//	services *services.Service
//}
//
//func NewHandler(services *services.Service) *Handler {
//	return &Handler{services: services}
//}

type Handler struct {
	services     *services.Service
	tokenManager pkg.TokenManager
}

func NewHandler(services *services.Service, tokenManager pkg.TokenManager) *Handler {
	return &Handler{
		services:     services,
		tokenManager: tokenManager,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.POST("/refresh", h.refresh)
	}

	api := router.Group("/api", h.userIdentity)
	{
		lists := api.Group("/lists")
		{
			lists.POST("/", h.createList)
			lists.GET("/", h.getAllLists)
			lists.GET("/:id", h.getListById)
			lists.PUT("/:id", h.updateList)
			lists.DELETE("/:id", h.deleteList)

			items := lists.Group(":id/items")
			{
				items.POST("/", h.createItem)
				items.GET("/", h.getAllItems)
			}
		}

		items := api.Group("items")
		{
			items.GET("/:id", h.getItemById)
			items.PUT("/:id", h.updateItem)
			items.DELETE("/:id", h.deleteItem)
		}
	}

	return router
}
