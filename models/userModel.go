package models

type User struct {
	Username        string `bson:"username"`
	Password        string `bson:"password"`
	Email           string `bson:"email"`
	FirstName       string `bson:"firstName"`
	LastName        string `bson:"lastName"`
	BirthYear       uint   `bson:"birthYear"`
	Oib             string `bson:"oib"`
	ActivationToken string `bson:"activationToken"`
	IsActivated     bool   `bson:"isActivated"`
}
