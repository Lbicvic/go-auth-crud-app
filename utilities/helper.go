package utilities

import (
	"log"
	"os"
	"time"

	"github.com/Lbicvic/go-auth-crud-app/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jinzhu/copier"
)

type ResponseUser struct {
	Username  string
	Email     string
	FirstName string
	LastName  string
	BirthYear uint
	Oib       string
}

func CreateToken(user_oib string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"oib":    user_oib,
		"expire": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("TOKEN_SECRET")))

	return tokenString, err
}

func GetResponseUserData(user *models.User) *ResponseUser {
	responseUser := &ResponseUser{}
	if err := copier.Copy(&responseUser, user); err != nil {
		log.Fatal(err)
	}
	return responseUser
}
