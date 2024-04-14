package controllers

import (
	"net/http"

	"github.com/Lbicvic/go-auth-crud-app/db/repositories"
	"github.com/Lbicvic/go-auth-crud-app/models"
	"github.com/Lbicvic/go-auth-crud-app/utilities"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	UserRepository    repositories.UserRepository
	ContactRepositroy repositories.ContactRepository
}

func ConstructUserController(userRepository repositories.UserRepository, contactRepository repositories.ContactRepository) UserController {
	return UserController{
		UserRepository:    userRepository,
		ContactRepositroy: contactRepository,
	}
}
func (userController *UserController) RegisterUser(context *gin.Context) {
	user := models.User{}
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := utilities.ValidateEmail(user.Email); err != nil {
		context.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	hashPass, er := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if er != nil {
		context.JSON(http.StatusBadGateway, gin.H{"message": er.Error()})
		return
	}
	user.Password = string(hashPass)
	user.ActivationToken = uuid.New().String()
	user.IsActivated = false
	err := userController.UserRepository.CreateUser(&user)
	if err != nil {
		context.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	tokenString, err := utilities.CreateToken(user.Oib)
	if err != nil {
		context.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	utilities.SendEmailVerification(user.Oib, user.FirstName, user.LastName, user.Email, user.ActivationToken)
	context.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func (userController *UserController) GetUser(context *gin.Context) {
	oib := context.Param("oib")
	user, err := userController.UserRepository.GetUserByOib(oib)
	if err != nil {
		context.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	responseUser := utilities.GetResponseUserData(user)
	context.JSON(http.StatusOK, responseUser)
}

func (userController *UserController) LoginUser(context *gin.Context) {
	var reqBody struct {
		Email    string
		Password string
	}
	if err := context.ShouldBindJSON(&reqBody); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	user, err := userController.UserRepository.GetUserByEmail(&reqBody.Email)
	if err != nil {
		context.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.Password)); err != nil {
		context.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	tokenString, err := utilities.CreateToken(user.Oib)
	if err != nil {
		context.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	responseUser := utilities.GetResponseUserData(user)
	context.JSON(http.StatusOK, gin.H{"user": responseUser, "token": tokenString})
}

func (userController *UserController) UpdateUser(context *gin.Context) {
	oib := context.Param("oib")
	user := models.User{}
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	updatedUser, err := userController.UserRepository.UpdateUser(&user, &oib)
	if err != nil {
		context.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	responseUser := utilities.GetResponseUserData(updatedUser)
	context.JSON(http.StatusOK, responseUser)
}

func (userController *UserController) DeleteUser(context *gin.Context) {
	oib := context.Param("oib")
	user, err := userController.UserRepository.GetUserByOib(oib)
	if err != nil {
		context.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	if user.IsActivated {
		utilities.SendEmailDeleteUser(user.Oib, user.FirstName, user.LastName, user.Email, user.ActivationToken)
		context.JSON(http.StatusOK, gin.H{"message": "Verification email has been sent for removing account"})
	} else {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "User is not verified, please verify account to proceed"})
	}
}

func (userController *UserController) ValidateEmailVerification(context *gin.Context) {
	oib := context.Param("oib")
	activationToken := context.Param("activationToken")
	user, err := userController.UserRepository.GetUserByOib(oib)
	if err != nil {
		context.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	if activationToken != user.ActivationToken {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "User not registered"})
		return
	}
	user.IsActivated = true
	_, er := userController.UserRepository.UpdateUser(user, &user.Oib)
	if er != nil {
		context.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "User successfully activated"})
}

func (userController *UserController) ValidateDeleteUser(context *gin.Context) {
	oib := context.Param("oib")
	activationToken := context.Param("activationToken")
	user, err := userController.UserRepository.GetUserByOib(oib)
	if err != nil {
		context.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	if user.IsActivated && user.ActivationToken == activationToken {
		if err := userController.ContactRepositroy.DeleteContacts(); err != nil {
			context.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
			return
		}
		err := userController.UserRepository.DeleteUser(&user.Oib)
		if err != nil {
			context.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
			return
		}
		context.JSON(http.StatusOK, gin.H{"message": "User successfully deleted"})
	} else {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "User is not verified, please verify account to proceed"})
	}
}
func (userController *UserController) PasswordRecovery(context *gin.Context) {
	email := context.Param("email")
	user, err := userController.UserRepository.GetUserByEmail(&email)
	if err != nil {
		context.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	utilities.SendEmailPassRecovery(user.Oib, user.FirstName, user.LastName, user.Email)
	context.JSON(http.StatusOK, gin.H{"message": "Email for password recovery has been sent"})
}

func (userController *UserController) ValidatePasswordRecovery(context *gin.Context) {
	oib := context.Param("oib")
	newPass := context.Param("newPass")
	user, err := userController.UserRepository.GetUserByOib(oib)
	if err != nil {
		context.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	if user.IsActivated {
		hashPass, er := bcrypt.GenerateFromPassword([]byte(newPass), 10)
		if er != nil {
			context.JSON(http.StatusBadGateway, gin.H{"message": er.Error()})
			return
		}
		user.Password = string(hashPass)
		_, err := userController.UserRepository.UpdateUser(user, &user.Oib)
		if err != nil {
			context.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
			return
		}
		context.JSON(http.StatusOK, gin.H{"message": "Password has been changed successfully"})
	} else {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "User is not verified, please verify account to proceed"})
	}
}
