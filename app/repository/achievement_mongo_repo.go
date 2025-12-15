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
	Create(detail *model.AchievementMongo) (primitive.ObjectID, error)
	Update(id primitive.ObjectID, detail map[string]interface{}) error
	FindByID(id primitive.ObjectID) (*model.AchievementMongo, error)
	AddAttachment(id primitive.ObjectID, att model.Attachment) error
}

type achievementMongoRepository struct {
	collection *mongo.Collection
}

func NewAchievementMongoRepository(db *mongo.Database) AchievementMongoRepository {
	return &achievementMongoRepository{
		collection: db.Collection("achievements"),
	}
}

// =====================
// CREATE DETAIL
// =====================
func (r *achievementMongoRepository) Create(detail *model.AchievementMongo) (primitive.ObjectID, error) {
	detail.CreatedAt = time.Now()
	detail.UpdatedAt = time.Now()

	res, err := r.collection.InsertOne(context.Background(), detail)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return res.InsertedID.(primitive.ObjectID), nil
}

// =====================
// UPDATE DETAIL
// =====================
func (r *achievementMongoRepository) Update(id primitive.ObjectID, detail map[string]interface{}) error {
	detail["updatedAt"] = time.Now()

	_, err := r.collection.UpdateByID(
		context.Background(),
		id,
		bson.M{"$set": detail},
	)
	return err
}

// =====================
// GET DETAIL
// =====================
func (r *achievementMongoRepository) FindByID(id primitive.ObjectID) (*model.AchievementMongo, error) {
	var result model.AchievementMongo
	err := r.collection.FindOne(
		context.Background(),
		bson.M{"_id": id},
	).Decode(&result)

	return &result, err
}

// =====================
// ADD ATTACHMENT
// =====================
func (r *achievementMongoRepository) AddAttachment(id primitive.ObjectID, att model.Attachment) error {
	att.UploadedAt = time.Now()

	_, err := r.collection.UpdateByID(
		context.Background(),
		id,
		bson.M{
			"$push": bson.M{"attachments": att},
			"$set":  bson.M{"updatedAt": time.Now()},
		},
	)
	return err
}
