package handler

import (
	"fmt"
	"github.com/Denialll/jwtauth-app/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) signUp(c *gin.Context) {
	var input models.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	id, err := h.services.Authorization.CreateUser(c.Request.Context(), input)
	fmt.Println("--------")
	fmt.Println(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

//type signInInput struct {
//	Username string `json:"username" binding:"required"`
//	Password string `json:"password" binding:"required"`
//}

func (h *Handler) signIn(c *gin.Context) {
	//var input signInInput
	//
	//if err := c.BindJSON(&input); err != nil {
	//	newErrorResponse(c, http.StatusBadRequest, err.Error())
	//	return
	//}

	uuid := c.Query("GUID")
	fmt.Println("guid: " + uuid)

	token, err := h.services.Authorization.GenerateTokens(c.Request.Context(), uuid)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}

func (h *Handler) refresh(c *gin.Context) {
}
