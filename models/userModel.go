package models

type User struct {
	Username  string `bson:"username"`
	Password  string `copier:"-" bson:"password"`
	Email     string `bson:"email"`
	FirstName string `bson:"firstName"`
	LastName  string `bson:"lastName"`
	BirthYear uint   `bson:"birthYear"`
	Oib       string `bson:"oib"`
}
