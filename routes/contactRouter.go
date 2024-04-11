package routes

import (
	"github.com/Lbicvic/go-auth-crud-app/db"
	"github.com/Lbicvic/go-auth-crud-app/middleware"
	"github.com/gin-gonic/gin"
)

func ContactRouter(apiRouter *gin.RouterGroup) {
	contactRoutes := apiRouter.Group("/contact")
	{
		contactRoutes.Use(middleware.RequireAuth)
		contactRoutes.POST("/create", db.ContactController.CreateContact)
		contactRoutes.GET("/:id", db.ContactController.GetContact)
		contactRoutes.POST("/", db.ContactController.GetContacts)
		contactRoutes.PATCH("/:id", db.ContactController.UpdateContact)
		contactRoutes.DELETE("/:id", db.ContactController.DeleteContact)
	}
}
