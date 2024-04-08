package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id        primitive.ObjectID `bson:"_id"`
	Username  string             `bson:"username"`
	Password  string             `bson:"password"`
	Email     string             `bson:"email"`
	FirstName string             `bson:"firstName"`
	LastName  string             `bson:"lastName"`
	BirthYear uint               `bson:"birthYear"`
	Oib       string             `bson:"oib"`
}
