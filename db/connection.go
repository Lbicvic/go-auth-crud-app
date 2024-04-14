package db

import (
	"context"
	"log"
	"os"

	"github.com/Lbicvic/go-auth-crud-app/controllers"
	"github.com/Lbicvic/go-auth-crud-app/db/repositories"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client            *mongo.Client
	UserRepository    repositories.UserRepository
	UserController    controllers.UserController
	contactRepository repositories.ContactRepository
	ContactController controllers.ContactController
)

func Connect() (*mongo.Collection, *mongo.Collection, context.Context) {
	var users *mongo.Collection
	var contacts *mongo.Collection
	ctx := context.TODO()
	Client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_ATLAS_URI")))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("DB connected")
	users = Client.Database("Phonebook").Collection("users")
	contacts = Client.Database("Phonebook").Collection("contacts")

	return users, contacts, ctx
}

func Initialize() {
	users, contacts, ctx := Connect()
	UserRepository = *repositories.ConstructUserRepository(users, ctx)
	contactRepository = *repositories.ConstructContactRepository(contacts, ctx)
	UserController = controllers.ConstructUserController(UserRepository, contactRepository)
	ContactController = controllers.ConstructContactController(contactRepository)
}
