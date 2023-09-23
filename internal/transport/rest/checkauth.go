package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary checkJWT
// @Security ApiKeyAuth
// @Tags checkJWT
// @Description JWT checker
// @ID checkJWT
// @Produce  json
// @Success 200 string text
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /checkjwt [get]
func (h *Handler) checkJWT(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"text": "Success!",
	})
}
