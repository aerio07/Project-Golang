package repository

import (
	"context"
	"time"

	"project_uas/app/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AchievementMongoRepository interface {
	Create(a *model.AchievementMongo) (primitive.ObjectID, error)
	FindByID(id primitive.ObjectID) (*model.AchievementMongo, error)
	Update(id primitive.ObjectID, update map[string]interface{}) error
	AddAttachment(id primitive.ObjectID, attachment model.Attachment) error
}

type achievementMongoRepository struct {
	collection *mongo.Collection
}

func NewAchievementMongoRepository(db *mongo.Database) AchievementMongoRepository {
	return &achievementMongoRepository{
		collection: db.Collection("achievements"),
	}
}

func (r *achievementMongoRepository) Create(a *model.AchievementMongo) (primitive.ObjectID, error) {
	now := time.Now()
	a.CreatedAt = now
	a.UpdatedAt = now
	a.Attachments = []model.Attachment{}

	res, err := r.collection.InsertOne(context.Background(), a)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return res.InsertedID.(primitive.ObjectID), nil
}

func (r *achievementMongoRepository) FindByID(id primitive.ObjectID) (*model.AchievementMongo, error) {
	var result model.AchievementMongo
	err := r.collection.FindOne(
		context.Background(),
		bson.M{"_id": id},
	).Decode(&result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *achievementMongoRepository) Update(id primitive.ObjectID, update map[string]interface{}) error {
	update["updatedAt"] = time.Now()

	_, err := r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": id},
		bson.M{
			"$set": update,
		},
	)

	return err
}

func (r *achievementMongoRepository) AddAttachment(id primitive.ObjectID, attachment model.Attachment) error {
	attachment.UploadedAt = time.Now()

	_, err := r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": id},
		bson.M{
			"$push": bson.M{
				"attachments": attachment,
			},
			"$set": bson.M{
				"updatedAt": time.Now(),
			},
		},
	)

	return err
}
