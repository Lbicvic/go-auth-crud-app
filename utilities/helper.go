package utilities

import (
	"errors"
	"log"
	"os"
	"time"

	emailverifier "github.com/AfterShip/email-verifier"
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

func ValidateEmail(email string) error {
	verifier := emailverifier.NewVerifier()
	result, err := verifier.Verify(email)
	if err != nil {
		return err
	}
	if !result.Syntax.Valid {
		return errors.New("email syntax invalid")
	}
	return nil
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
	from := mail.NewEmail("No Reply", os.Getenv("SENDGRID_SENDER"))
	subject := "Sign up email verification"
	to := mail.NewEmail(firstName+lastName, email)
	htmlContent := "<h1>Click link to verify your email</h1> <a href=" + `http://localhost:` + domainPort + `/api/auth/activate/` + oib + `/` + activationToken + ">Verify Email</a>"
	message := mail.NewSingleEmail(from, subject, to, "", htmlContent)
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

func SendEmailPassRecovery(oib string, firstName string, lastName string, email string) {
	domainPort := os.Getenv("PORT")
	from := mail.NewEmail("No Reply", os.Getenv("SENDGRID_SENDER"))
	subject := "Password Recovery"
	to := mail.NewEmail(firstName+lastName, email)
	htmlContent := "<h1>Click link to change your password</h1> <a href=" + `http://localhost:` + domainPort + `/api/user/forgotPass/` + oib + `/` + ">Change Password</a>"
	message := mail.NewSingleEmail(from, subject, to, "", htmlContent)
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

func SendEmailDeleteUser(oib string, firstName string, lastName string, email string, token string) {
	domainPort := os.Getenv("PORT")
	from := mail.NewEmail("No Reply", os.Getenv("SENDGRID_SENDER"))
	subject := "Remove Account Verification"
	to := mail.NewEmail(firstName+lastName, email)
	htmlContent := "<h1>Click link to remove your account</h1> <a href=" + `http://localhost:` + domainPort + `/api/user/authorizeDelete/` + oib + `/` + token + ">Confirm</a>"
	message := mail.NewSingleEmail(from, subject, to, "", htmlContent)
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
