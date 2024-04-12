package utilities

import (
	"log"
	"os"
	"time"

	"github.com/Lbicvic/go-auth-crud-app/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jinzhu/copier"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
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

func SendEmailVerification(oib string, firstName string, lastName string, email string, activationToken string) {
	domainPort := os.Getenv("PORT")
	from := mail.NewEmail("No Reply", "leopold.bicvic@gmail.com")
	subject := "Sign up email verification"
	to := mail.NewEmail(firstName+lastName, email)
	plainTextContent := "and easy to do anywhere, even with Go"
	htmlContent := "<h1>Click link to verify your email</h1> <a href=" + `http://localhost:` + domainPort + `/api/auth/activate/` + oib + `/` + activationToken + ">Verify Email</a>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		log.Println(response.StatusCode)
		log.Println(response.Body)
		log.Println(response.Headers)
	}
}
