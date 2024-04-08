package db

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_ATLAS_URI")))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("DB connected")

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}
