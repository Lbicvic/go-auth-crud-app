package routes

import (
	"github.com/Lbicvic/go-auth-crud-app/db"
	"github.com/Lbicvic/go-auth-crud-app/middleware"
	"github.com/gin-gonic/gin"
)

func UserRouter(apiRouter *gin.RouterGroup) {
	userRoutes := apiRouter.Group("/user")
	{
		userRoutes.Use(middleware.RequireAuth)
		userRoutes.GET("/:id", db.UserController.GetUser)
		userRoutes.PATCH("/:id", db.UserController.UpdateUser)
		userRoutes.DELETE("/:id", db.UserController.DeleteUser)
	}
}
