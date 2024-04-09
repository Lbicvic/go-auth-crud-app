package repositories

import (
	"context"

	"github.com/Lbicvic/go-auth-crud-app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IUserRepository interface {
	CreateUser(*models.User) error
	GetUserById(*string) (*models.User, error)
	GetUserByEmail(*string) (*models.User, error)
	UpdateUser(*models.User) error
	DeleteUser(*string) error
}
type UserRepository struct {
	users   *mongo.Collection
	context context.Context
}

func ConstructUserRepository(users *mongo.Collection, context context.Context) *UserRepository {
	return &UserRepository{
		users:   users,
		context: context,
	}
}

func (userRepository *UserRepository) CreateUser(user *models.User) error {
	_, err := userRepository.users.InsertOne(userRepository.context, user)
	return err
}

func (userRepository *UserRepository) GetUserById(id *string) (*models.User, error) {
	var user *models.User
	userId, err := primitive.ObjectIDFromHex(*id)
	if err != nil {
		return user, err
	}
	filter := bson.D{bson.E{Key: "_id", Value: userId}}
	err = userRepository.users.FindOne(userRepository.context, filter).Decode(&user)
	return user, err
}

func (userRepository *UserRepository) GetUserByEmail(email *string) (*models.User, error) {
	var user *models.User
	filter := bson.D{bson.E{Key: "email", Value: email}}
	err := userRepository.users.FindOne(userRepository.context, filter).Decode(&user)
	return user, err
}

func (userRepository *UserRepository) UpdateUser(user *models.User) (*models.User, error) {
	return nil, nil
}

func (userRepository *UserRepository) DeleteUser(id *string) error {
	return nil
}
