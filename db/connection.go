package db

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

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
