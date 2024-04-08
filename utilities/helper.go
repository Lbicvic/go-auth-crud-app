package utilities

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(user_oib string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"oib":    user_oib,
		"expire": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(os.Getenv("TOKEN_SECRET"))

	return tokenString, err
}
