package mongodb

import (
	"context"
	"devconnectstorage/internal/domain"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoFileRepository struct {
	client     *mongo.Client
	database   string
	collection string
}

func NewMongoFileRepository(
	mongoUri string,
	username string,
	password string,
	database string,
	collection string,
) (*MongoFileRepository, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoUri))
	if err != nil {
		return &MongoFileRepository{}, err
	}

	return &MongoFileRepository{
		client:     client,
		database:   database,
		collection: collection,
	}, nil
}

func (repo MongoFileRepository) Save(ctx context.Context, file domain.File) (domain.File, error) {
	result, err := repo.client.Database(repo.database).Collection(repo.collection).InsertOne(ctx, NewMongoFileEntity(file))
	if err != nil {
		return domain.File{}, err
	}

	if result.InsertedID == nil {
		return domain.File{}, errors.New("failed to insert file")
	}
	return file, nil
}
