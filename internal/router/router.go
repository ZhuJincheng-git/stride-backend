package router

import (
	"github.com/ZhuJincheng-git/stride-backend/internal/handler"
	"github.com/ZhuJincheng-git/stride-backend/internal/middleware"
	"github.com/ZhuJincheng-git/stride-backend/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Auth *handler.AuthHandler
}

func Build(h *Handlers, tokens *jwt.Manager) *gin.Engine {
	r := gin.New()

	// Global Middleware
	r.Use(gin.Logger(), gin.Recovery())

	// Liveness probe
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API Version 1
	v1 := r.Group("/api/v1")

	// --- public routes ---
	auth := v1.Group("/auth")
	{
		auth.POST("/register", h.Auth.Register)
		auth.POST("/login", h.Auth.Login)
	}

	// --- authenticated routes ---
	authed := v1.Group("")
	authed.Use(middleware.AuthRequired(tokens))
	{
		authed.GET("/auth/me", h.Auth.Me)
	}

	return r
}