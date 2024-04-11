package controllers

import (
	"net/http"

	"github.com/Lbicvic/go-auth-crud-app/db/repositories"
	"github.com/Lbicvic/go-auth-crud-app/models"
	"github.com/gin-gonic/gin"
)

type ContactController struct {
	ContactRepository repositories.ContactRepository
}

func ConstructContactController(contactRepository repositories.ContactRepository) ContactController {
	return ContactController{
		ContactRepository: contactRepository,
	}
}
func (contactController *ContactController) CreateContact(context *gin.Context) {
	contact := models.Contact{}
	if err := context.ShouldBindJSON(&contact); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err := contactController.ContactRepository.CreateContact(&contact)
	if err != nil {
		context.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Contact has been created"})
}

func (contactController *ContactController) GetContact(context *gin.Context) {
	_id := context.Param("id")
	contact, err := contactController.ContactRepository.GetContact(&_id)
	if err != nil {
		context.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	context.JSON(http.StatusOK, contact)
}

func (contactController *ContactController) GetContacts(context *gin.Context) {
	var reqBody struct {
		User_oib string
	}
	if err := context.ShouldBindJSON(&reqBody); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	contacts, err := contactController.ContactRepository.GetContacts(&reqBody.User_oib)
	if err != nil {
		context.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	context.JSON(http.StatusOK, contacts)
}
func (contactController *ContactController) UpdateContact(context *gin.Context) {
	_id := context.Param("id")
	contact := models.Contact{}
	if err := context.ShouldBindJSON(&contact); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	updatedContact, err := contactController.ContactRepository.UpdateContact(&contact, &_id)
	if err != nil {
		context.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	context.JSON(http.StatusOK, updatedContact)
}

func (contactController *ContactController) DeleteContact(context *gin.Context) {
	_id := context.Param("id")
	err := contactController.ContactRepository.DeleteContact(&_id)
	if err != nil {
		context.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Contact successfully deleted"})
}
