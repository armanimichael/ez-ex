package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type accountRoutes struct {
	authService Authenticator
}

type hashPasswordBody struct {
	Password string `json:"password" binding:"required"`
}

func newAccountRoutes(handler *gin.RouterGroup, authService Authenticator) {
	routes := accountRoutes{
		authService: authService,
	}

	handler.POST("/hash-password", routes.hashAndSaltPassword)
}

func (r *accountRoutes) hashAndSaltPassword(c *gin.Context) {
	var json hashPasswordBody
	if err := c.BindJSON(&json); err != nil {
		return
	}

	encode := r.authService.HashAndSaltPassword(c.Request.Context(), json.Password)
	c.JSON(http.StatusOK, gin.H{
		"result": encode,
	})
}
