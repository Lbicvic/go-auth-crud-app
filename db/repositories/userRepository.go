package repositories

import (
	"context"
	"log"

	"github.com/Lbicvic/go-auth-crud-app/models"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IUserRepository interface {
	CreateUser(*models.User) error
	GetUserById(*string) (*models.User, error)
	GetUserByEmail(*string) (*models.User, error)
	UpdateUser(*models.User, *string) error
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

func (userRepository *UserRepository) UpdateUser(user *models.User, id *string) (*models.User, error) {
	var updatedUser *models.User
	userId, err := primitive.ObjectIDFromHex(*id)
	if err != nil {
		log.Fatal(err)
	}
	currentUser, err := userRepository.GetUserById(id)
	if err != nil {
		log.Fatal(err)
	}
	if err := copier.CopyWithOption(currentUser, user, copier.Option{IgnoreEmpty: true}); err != nil {
		log.Fatal(err)
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	filter := bson.D{bson.E{Key: "_id", Value: userId}}
	update := bson.D{bson.E{Key: "$set", Value: currentUser}}
	err = userRepository.users.FindOneAndUpdate(userRepository.context, filter, update, opts).Decode(&updatedUser)
	if err != nil {
		log.Fatal(err)
	}
	return updatedUser, err
}

func (userRepository *UserRepository) DeleteUser(id *string) error {
	userId, err := primitive.ObjectIDFromHex(*id)
	if err != nil {
		return err
	}
	filter := bson.D{bson.E{Key: "_id", Value: userId}}
	result, err := userRepository.users.DeleteOne(userRepository.context, filter)
	if result.DeletedCount != 1 {
		return err
	}
	return nil
}
