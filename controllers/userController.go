package controllers

import (
	"net/http"

	"github.com/Lbicvic/go-auth-crud-app/db/repositories"
	"github.com/Lbicvic/go-auth-crud-app/models"
	"github.com/Lbicvic/go-auth-crud-app/utilities"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	UserRepository repositories.UserRepository
}

func constructUserController(userRepository repositories.UserRepository) UserController {
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
	context.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func (userController *UserController) GetUser(context *gin.Context) {
	_id := context.Param("id")
	user, err := userController.UserRepository.GetUser(&_id)
	if err != nil {
		context.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	context.JSON(http.StatusOK, user)
}

func (userController *UserController) LoginUser(context *gin.Context) {

}

func (userController *UserController) UpdateUser(context *gin.Context) {

}

func (userController *UserController) DeleteUser(context *gin.Context) {

}

func (userController *UserController) UserRoutes(apiRouter *gin.RouterGroup) {
	userRoutes := apiRouter.Group("/user")
	{
		userRoutes.POST("/register", userController.RegisterUser)
		userRoutes.POST("/login", userController.LoginUser)
		userRoutes.GET("/:id", userController.GetUser)
		userRoutes.PATCH("/:id", userController.UpdateUser)
		userRoutes.DELETE("/login", userController.DeleteUser)
	}
}
