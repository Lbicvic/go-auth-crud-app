package main

import (
	"log"
	"os"

	"github.com/Lbicvic/go-auth-crud-app/db"
	"github.com/Lbicvic/go-auth-crud-app/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db.Connect()
}
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router := gin.Default()

	routes.UserRoutes(router)

	router.Run(":" + port)
}
