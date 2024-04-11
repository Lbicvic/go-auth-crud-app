package main

import (
	"context"
	"log"
	"os"

	"github.com/Lbicvic/go-auth-crud-app/db"
	"github.com/Lbicvic/go-auth-crud-app/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	ctx context.Context
	err error
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db.Initialize()
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
	routes.AuthRouter(apiPath)
	routes.UserRouter(apiPath)
	routes.ContactRouter(apiPath)
	router.Run(":" + port)
}
