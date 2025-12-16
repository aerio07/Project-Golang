package repository

import (
	"context"

	"project_uas/app/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReportMongoRepository interface {
	AggregateStatistics(studentIDs []string) (*model.AchievementStatistics, error)
}

type reportMongoRepository struct {
	collection *mongo.Collection
}

func NewReportMongoRepository(db *mongo.Database) ReportMongoRepository {
	return &reportMongoRepository{
		collection: db.Collection("achievements"),
	}
}

func (r *reportMongoRepository) AggregateStatistics(studentIDs []string) (*model.AchievementStatistics, error) {
	match := bson.D{}
	if len(studentIDs) > 0 {
		match = bson.D{{Key: "studentId", Value: bson.M{"$in": studentIDs}}}
	}

	// pakai $facet biar 1 kali hit
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: match}},
		{{
			Key: "$facet",
			Value: bson.M{
				"totalPerType": mongo.Pipeline{
					{{Key: "$group", Value: bson.M{"_id": "$achievementType", "total": bson.M{"$sum": 1}}}},
					{{Key: "$sort", Value: bson.M{"total": -1}}},
				},
				"totalPerPeriod": mongo.Pipeline{
					{{Key: "$group", Value: bson.M{
						"_id": bson.M{
							"$dateToString": bson.M{"format": "%Y-%m", "date": "$createdAt"},
						},
						"total": bson.M{"$sum": 1},
					}}},
					{{Key: "$sort", Value: bson.M{"_id": 1}}},
				},
				"topStudents": mongo.Pipeline{
					{{Key: "$group", Value: bson.M{
						"_id":    "$studentId",
						"points": bson.M{"$sum": "$points"},
						"total":  bson.M{"$sum": 1},
					}}},
					{{Key: "$sort", Value: bson.M{"points": -1}}},
					{{Key: "$limit", Value: 5}},
				},
				"levelDist": mongo.Pipeline{
					// asumsi tingkat kompetisi disimpan di details.level
					{{Key: "$group", Value: bson.M{"_id": "$details.level", "total": bson.M{"$sum": 1}}}},
					{{Key: "$sort", Value: bson.M{"total": -1}}},
				},
			},
		}},
	}

	cur, err := r.collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	var res []struct {
		TotalPerType   []model.CountByKey `bson:"totalPerType"`
		TotalPerPeriod []model.CountByKey `bson:"totalPerPeriod"`
		TopStudents    []model.TopStudent `bson:"topStudents"`
		LevelDist      []model.CountByKey `bson:"levelDist"`
	}

	if err := cur.All(context.Background(), &res); err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return &model.AchievementStatistics{}, nil
	}

	out := &model.AchievementStatistics{
		TotalPerType:   res[0].TotalPerType,
		TotalPerPeriod: res[0].TotalPerPeriod,
		TopStudents:    res[0].TopStudents,
		LevelDist:      res[0].LevelDist,
	}
	return out, nil
}
