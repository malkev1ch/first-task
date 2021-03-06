package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/malkev1ch/first-task/internal/model"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

// CatRepositoryMongo type represents mongo object cat structure and behavior.
type CatRepositoryMongo struct {
	DB *mongo.Client
}

func NewCatRepositoryMongo(db *mongo.Client) *CatRepositoryMongo {
	return &CatRepositoryMongo{
		DB: db,
	}
}

// Create method saves object Cat into mongo database.
func (r CatRepositoryMongo) Create(ctx context.Context, input *model.Cat) error {
	logrus.WithFields(logrus.Fields{
		"Name":       input.Name,
		"DateBirth":  input.DateBirth,
		"Vaccinated": input.Vaccinated,
	}).Debugf("mongo repository: create cat")
	col := r.DB.Database("mongo_database").Collection("cats")

	_, err := col.InsertOne(ctx, bson.D{
		{Key: "_id", Value: input.ID},
		{Key: "name", Value: input.Name},
		{Key: "dateBirth", Value: input.DateBirth},
		{Key: "vaccinated", Value: input.Vaccinated},
	})
	if err != nil {
		logrus.Error(err, "mongo repository: Error occurred while inserting new row in table cats")
		return fmt.Errorf("mongo repository: can't create cat - %w", err)
	}

	return nil
}

// Get method returns object Cat from mongo database
// with selection by id.
func (r CatRepositoryMongo) Get(ctx context.Context, id string) (*model.Cat, error) {
	logrus.WithFields(logrus.Fields{
		"ID": id,
	}).Debugf("mongo repository: get cat")
	col := r.DB.Database("mongo_database").Collection("cats")
	var cat model.Cat
	err := col.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&cat)
	if err != nil {
		logrus.Error(err, "mongo repository: Error occurred while selecting row from table cats")
		return nil, fmt.Errorf("mongo repository: can't get cat - %w", err)
	}
	return &cat, nil
}

// Update method updates object Cat from mongo database
// with selection by id.
func (r CatRepositoryMongo) Update(ctx context.Context, id string, input *model.UpdateCat) (*model.Cat, error) {
	logrus.WithFields(logrus.Fields{
		"Name":       input.Name,
		"DateBirth":  input.DateBirth,
		"Vaccinated": input.Vaccinated,
	}).Debugf("mongo repository: update cat")

	col := r.DB.Database("mongo_database").Collection("cats")

	_, err := col.UpdateOne(ctx, bson.D{{Key: "_id", Value: id}}, bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "name", Value: input.Name},
			{Key: "dateBirth", Value: input.DateBirth},
			{Key: "vaccinated", Value: input.Vaccinated},
		}},
	})
	if err != nil {
		logrus.Error(err, "mongo repository: Error occurred while updating row from table cats")
		return nil, fmt.Errorf("mongo repository: can't update cat - %w", err)
	}
	var cat model.Cat
	if err := col.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&cat); err != nil {
		logrus.Error(err, "mongo repository: Error occurred while selecting row from table cats")
		return nil, fmt.Errorf("mongo repository: can't get cat - %w", err)
	}

	return &cat, nil
}

// Delete method deletes object Cat from mongo database
// with selection by id.
func (r CatRepositoryMongo) Delete(ctx context.Context, id string) error {
	logrus.WithFields(logrus.Fields{
		"ID": id,
	}).Debugf("mongo repository: delete cat")
	col := r.DB.Database("mongo_database").Collection("cats")
	_, err := col.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})
	if err != nil {
		logrus.Error(err, "Error occurred while deleting row from table cats")
		return fmt.Errorf("mongodb repository: can't delete cat - %w", err)
	}
	return nil
}

// UploadImage method updates image path object Cat from mongo database
// with selection by id.
func (r CatRepositoryMongo) UploadImage(ctx context.Context, id string, path string) (*model.Cat, error) {
	logrus.WithFields(logrus.Fields{
		"ID": id,
	}).Debugf("mongo repository: update cats image path")
	col := r.DB.Database("mongo_database").Collection("cats")
	_, err := col.UpdateOne(ctx, bson.D{{Key: "_id", Value: id}}, bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "imagePath", Value: path},
		}},
	})
	if err != nil {
		logrus.Error(err, "mongo repository: Error occurred while updating image path table cats")
		return nil, fmt.Errorf("mongo repository: can't update cats image path - %w", err)
	}
	var cat model.Cat
	if err := col.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&cat); err != nil {
		logrus.Error(err, "mongo repository: Error occurred while selecting row from table cats")
		return nil, fmt.Errorf("mongo repository: can't get cat - %w", err)
	}

	return &cat, nil
}
