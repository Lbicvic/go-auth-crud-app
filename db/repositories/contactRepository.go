package repositories

import "github.com/Lbicvic/go-auth-crud-app/models"

type ContactRepository interface {
	CreateContact(*models.Contact) error
	GetContacts(*string) ([]*models.Contact, error)
	GetContact(*string) (*models.Contact, error)
	UpdateContact(*models.Contact) error
	DeleteContact(*string) error
}
