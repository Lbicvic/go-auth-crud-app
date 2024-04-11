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

type IContactRepository interface {
	CreateContact(*models.Contact) error
	GetContacts(*string) ([]*models.Contact, error)
	GetContact(*string) (*models.Contact, error)
	UpdateContact(*models.Contact, *string) error
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
	opts := options.Find().SetSort(bson.D{bson.E{Key: "firstName", Value: 1}})
	cursor, err := contactRepository.contacts.Find(contactRepository.context, filter, opts)
	if err := cursor.All(contactRepository.context, &contacts); err != nil {
		log.Fatal(err)
	}
	return contacts, err
}

func (contactRepository *ContactRepository) UpdateContact(contact *models.Contact, id *string) (*models.Contact, error) {
	var updatedContact *models.Contact
	contactId, err := primitive.ObjectIDFromHex(*id)
	if err != nil {
		return contact, err
	}
	currentContact, err := contactRepository.GetContact(id)
	if err != nil {
		return nil, err
	}
	if err := copier.CopyWithOption(currentContact, contact, copier.Option{IgnoreEmpty: true}); err != nil {
		log.Fatal(err)
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	filter := bson.D{bson.E{Key: "_id", Value: contactId}}
	update := bson.D{bson.E{Key: "$set", Value: currentContact}}
	err = contactRepository.contacts.FindOneAndUpdate(contactRepository.context, filter, update, opts).Decode(&updatedContact)
	if err != nil {
		log.Fatal(err)
	}
	return updatedContact, err
}

func (contactRepository *ContactRepository) DeleteContact(id *string) error {
	contactId, err := primitive.ObjectIDFromHex(*id)
	if err != nil {
		return err
	}
	filter := bson.D{bson.E{Key: "_id", Value: contactId}}
	result, err := contactRepository.contacts.DeleteOne(contactRepository.context, filter)
	if result.DeletedCount != 1 {
		return err
	}
	return nil
}
