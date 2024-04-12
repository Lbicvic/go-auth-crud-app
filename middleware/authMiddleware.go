package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Lbicvic/go-auth-crud-app/db"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(context *gin.Context) {
	authorization := context.GetHeader("Authorization")
	if authorization == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "User not authorized, please log in to proceed"})
	}

	tokenString := strings.Split(authorization, " ")[1]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})
	if err != nil {
		log.Fatal(err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if float64(time.Now().Unix()) > claims["expire"].(float64) {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "token has expired, please log in again"})
		}
		user, err := db.UserRepository.GetUserByOib(claims["oib"].(string))
		if err != nil {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "User not authorized, please log in to proceed"})
		}
		if !user.IsActivated {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "User not verified, please verfy your email adress to proceed"})
		}
		context.Next()
	} else {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "User not authorized, please log in to proceed"})
	}

}
