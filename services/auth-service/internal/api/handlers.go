package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/gmsas95/blytz-mvp/services/auth-service/internal/models"
	"github.com/gmsas95/blytz-mvp/services/auth-service/internal/services"
	shared_errors "github.com/gmsas95/blytz-mvp/shared/pkg/errors"
	"github.com/gmsas95/blytz-mvp/shared/pkg/utils"
)

// RegisterRoutes registers the API routes for the auth service
func RegisterRoutes(router *gin.Engine, authService *services.AuthService) {
	api := router.Group("/api/v1")
	{
		api.POST("/register", registerUser(authService))
		api.POST("/login", loginUser(authService))
	}
}

func registerUser(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			utils.SendErrorResponse(c, shared_errors.ErrInvalidRequestBody)
			return
		}

		if err := authService.RegisterUser(&user); err != nil {
			utils.SendErrorResponse(c, err)
			return
		}

		utils.SendSuccessResponse(c, http.StatusCreated, gin.H{"message": "User registered successfully"})
	}
}

func loginUser(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginDetails struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&loginDetails); err != nil {
			utils.SendErrorResponse(c, shared_errors.ErrInvalidRequestBody)
			return
		}

		token, err := authService.LoginUser(loginDetails.Email, loginDetails.Password)
		if err != nil {
			utils.SendErrorResponse(c, err)
			return
		}

		utils.SendSuccessResponse(c, http.StatusOK, gin.H{"token": token})
	}
}

