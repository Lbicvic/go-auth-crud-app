package repositories

import (
	"context"
	"log"

	"github.com/Lbicvic/go-auth-crud-app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IContactRepository interface {
	CreateContact(*models.Contact) error
	GetContacts(*string) ([]*models.Contact, error)
	GetContact(*string) (*models.Contact, error)
	UpdateContact(*models.Contact) error
	DeleteContact(*string) error
}
type ContactRepository struct {
	contacts *mongo.Collection
	context  context.Context
}

func ConstructContactRepository(contacts *mongo.Collection, context context.Context) *ContactRepository {
	return &ContactRepository{
		contacts: contacts,
		context:  context,
	}
}

func (contactRepository *ContactRepository) CreateContact(contact *models.Contact) error {
	_, err := contactRepository.contacts.InsertOne(contactRepository.context, contact)
	return err
}

func (contactRepository *ContactRepository) GetContact(id *string) (*models.Contact, error) {
	var contact *models.Contact
	contactId, err := primitive.ObjectIDFromHex(*id)
	if err != nil {
		return contact, err
	}
	filter := bson.D{bson.E{Key: "_id", Value: contactId}}
	err = contactRepository.contacts.FindOne(contactRepository.context, filter).Decode(&contact)
	return contact, err
}

func (contactRepository *ContactRepository) GetContacts(user_oib *string) ([]*models.Contact, error) {
	var contacts []*models.Contact
	filter := bson.D{{Key: "user_oib", Value: user_oib}}
	opts := options.Find().SetSort(bson.D{{"firstName", 1}})
	cursor, err := contactRepository.contacts.Find(contactRepository.context, filter, opts)
	if err := cursor.All(contactRepository.context, &contacts); err != nil {
		log.Fatal(err)
	}
	return contacts, err
}

func (contactRepository *ContactRepository) UpdateContact(contact *models.Contact) (*models.Contact, error) {
	return nil, nil
}

func (contactRepository *ContactRepository) DeleteContact(id *string) error {
	return nil
}
