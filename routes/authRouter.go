package routes

import (
	"github.com/Lbicvic/go-auth-crud-app/db"
	"github.com/gin-gonic/gin"
)

func AuthRouter(apiRouter *gin.RouterGroup) {
	authRoutes := apiRouter.Group("/auth")
	{
		authRoutes.POST("/register", db.UserController.RegisterUser)
		authRoutes.POST("/login", db.UserController.LoginUser)
		authRoutes.GET("/activate/:oib/:activationToken", db.UserController.ValidateEmailVerification)
	}
}
