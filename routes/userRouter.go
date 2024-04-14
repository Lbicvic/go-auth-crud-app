package routes

import (
	"github.com/Lbicvic/go-auth-crud-app/db"
	"github.com/Lbicvic/go-auth-crud-app/middleware"
	"github.com/gin-gonic/gin"
)

func UserRouter(apiRouter *gin.RouterGroup) {
	userRoutes := apiRouter.Group("/user")
	{
		userRoutes.GET("/authorizeDelete/:oib/:activationToken", db.UserController.ValidateDeleteUser)
		userRoutes.Use(middleware.RequireAuth)
		userRoutes.GET("/:oib", db.UserController.GetUser)
		userRoutes.PATCH("/:oib", db.UserController.UpdateUser)
		userRoutes.DELETE("/:oib", db.UserController.DeleteUser)
		userRoutes.GET("/passRecovery/:email", db.UserController.PasswordRecovery)
		userRoutes.POST("/passRecovery/:oib/:newPass", db.UserController.ValidatePasswordRecovery)
	}
}
