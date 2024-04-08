package repositories

import "github.com/Lbicvic/go-auth-crud-app/models"

type UserRepository interface {
	CreateUser(*models.User) error
	GetUser(string) (*models.User, error)
	UpdateUser(*models.User) error
	DeleteUser(string) error
}
