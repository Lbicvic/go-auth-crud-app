package models

type Contact struct {
	FirstName   string `bson:"firstName"`
	LastName    string `bson:"lastName"`
	PhoneNumber string `bson:"phoneNumber"`
	Email       string `bson:"email"`
	User_oib    string `bson:"user_oib"`
}
