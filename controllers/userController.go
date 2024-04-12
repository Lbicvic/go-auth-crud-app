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
	UserRepository repositories.UserRepository
}

func ConstructUserController(userRepository repositories.UserRepository) UserController {
	return UserController{
		UserRepository: userRepository,
	}
}
func (userController *UserController) RegisterUser(context *gin.Context) {
	user := models.User{}
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	hashPass, er := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if er != nil {
		context.JSON(http.StatusBadGateway, gin.H{"message": er.Error()})
		return
	}
	user.Password = string(hashPass)
	hashActivationToken, er := bcrypt.GenerateFromPassword([]byte(uuid.New().String()), 10)
	if er != nil {
		context.JSON(http.StatusBadGateway, gin.H{"message": er.Error()})
		return
	}
	user.ActivationToken = string(hashActivationToken)
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
	context.JSON(http.StatusOK, gin.H{"token": tokenString, "activationToken": user.ActivationToken})
}

func (userController *UserController) GetUser(context *gin.Context) {
	_id := context.Param("id")
	user, err := userController.UserRepository.GetUserById(&_id)
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
	_id := context.Param("id")
	user := models.User{}
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	updatedUser, err := userController.UserRepository.UpdateUser(&user, &_id)
	if err != nil {
		context.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	responseUser := utilities.GetResponseUserData(updatedUser)
	context.JSON(http.StatusOK, responseUser)
}

func (userController *UserController) DeleteUser(context *gin.Context) {
	_id := context.Param("id")
	err := userController.UserRepository.DeleteUser(&_id)
	if err != nil {
		context.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "User successfully deleted"})
}

func (userController *UserController) ValidateEmailVerification(context *gin.Context) {
	oib := context.Param("oib")
	activationToken := context.Param("activationToken")
	user, err := userController.UserRepository.GetUserByOib(oib)
	if err != nil {
		context.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.ActivationToken), []byte(activationToken)); err != nil {
		context.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	user.IsActivated = true
	context.JSON(http.StatusOK, gin.H{"message": "User successfully activated"})
}
