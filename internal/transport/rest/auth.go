package handler

import (
	"github.com/Denialll/jwtauth-app/internal/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type signUpInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary SignUp
// @Tags auth
// @Description User create
// @ID create-account
// @Accept  json
// @Produce  json
// @Param input body signUpInput true "User Info"
// @Success 200 {string} string id
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var input model.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	id, err := h.services.Authorization.CreateUser(c.Request.Context(), input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary SignIn
// @Tags auth
// @Description User auth
// @ID user-auth
// @Produce  json
// @Param guid query string true "User GUID"
// @Success 200
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	guid := c.Query("guid")

	token, err := h.services.Authorization.GenerateTokens(c.Request.Context(), guid)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}

type refreshInput struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// @Summary Refresh
// Security ApiKeyAuth
// @Tags auth
// @Description RefreshTokens
// @ID refreshTokens
// @Accept  json
// @Produce  json
// @Param Authorization header string true "AccessToken"
// @Param input body refreshInput true "User Info"
// @Success 200
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/refresh [post]
func (h *Handler) refresh(c *gin.Context) {
	var input refreshInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.services.Authorization.UpdateTokens(c.Request.Context(), c.GetHeader("Authorization"), input.RefreshToken)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
