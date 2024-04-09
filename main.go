package main

import (
	"context"
	"log"
	"os"

	"github.com/Lbicvic/go-auth-crud-app/controllers"
	"github.com/Lbicvic/go-auth-crud-app/db"
	"github.com/Lbicvic/go-auth-crud-app/db/repositories"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	ctx               context.Context
	userRepository    repositories.UserRepository
	userController    controllers.UserController
	contactRepository repositories.ContactRepository
	contactController controllers.ContactController
	err               error
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	users, contacts, ctx := db.Connect()
	userRepository = *repositories.ConstructUserRepository(users, ctx)
	userController = controllers.ConstructUserController(userRepository)
	contactRepository = *repositories.ConstructContactRepository(contacts, ctx)
	contactController = controllers.ConstructContactController(contactRepository)
}
func main() {
	defer func() {
		if err = db.Client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router := gin.Default()
	apiPath := router.Group("/api")
	userController.UserRoutes(apiPath)
	contactController.ContactRoutes(apiPath)
	router.Run(":" + port)
}
