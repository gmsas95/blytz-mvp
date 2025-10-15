package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/blytz/auth-service/internal/api/handlers"
	"github.com/blytz/auth-service/internal/middleware"
)

func SetupRoutes(router *gin.Engine, authHandler *handlers.AuthHandler) {
	// Public routes
	public := router.Group("/api/v1/auth")
	{
		public.POST("/signup", authHandler.SignUp)
		public.POST("/login", authHandler.Login)
		public.POST("/refresh", authHandler.RefreshToken)
	}

	// Protected routes
	protected := router.Group("/api/v1/auth")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/verify", authHandler.Verify)
		protected.POST("/logout", authHandler.Logout)
		protected.PUT("/profile", authHandler.UpdateProfile)
		protected.GET("/profile", authHandler.GetProfile)
	}
}
