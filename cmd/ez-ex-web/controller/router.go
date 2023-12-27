package controller

import (
	"context"
	"github.com/gin-gonic/gin"
)

type Authenticator interface {
	HashAndSaltPassword(ctx context.Context, password string) string
}

func New(authService Authenticator) *gin.Engine {
	engine := gin.New()
	g := engine.Group("/api")

	// Account related requests
	newAccountRoutes(g.Group("/account"), authService)

	return engine
}
