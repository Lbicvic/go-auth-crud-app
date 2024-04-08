package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Contact struct {
	Id          primitive.ObjectID `bson:"_id"`
	FirstName   string             `bson:"firstName"`
	LastName    string             `bson:"lastName"`
	PhoneNumber string             `bson:"phoneNumber"`
	Email       string             `bson:"email"`
}
