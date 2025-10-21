package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/gmsas95/blytz-mvp/services/auth-service/internal/api/handlers"
	"github.com/gmsas95/blytz-mvp/services/auth-service/internal/middleware"
	"github.com/gmsas95/blytz-mvp/services/auth-service/internal/services"
)

func SetupRoutes(router *gin.Engine, authHandler *handlers.AuthHandler, authService *services.AuthService) {
	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		protected := api.Group("/protected").Use(middleware.AuthMiddleware(authService))
		{
			protected.GET("/profile", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "this is a protected route"})
			})
		}
	}
}
