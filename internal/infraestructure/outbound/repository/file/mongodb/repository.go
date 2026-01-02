package mongodb

import (
	"context"
	"devconnectstorage/internal/domain"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
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

func (repo MongoFileRepository) GetFile(ctx context.Context, id string) (domain.File, error) {
	filter := bson.M{"_id": id}
	result := repo.client.Database(repo.database).Collection(repo.collection).FindOne(ctx, filter)
	if result.Err() != nil {
		return domain.File{}, result.Err()
	}

	var mongoFile MongoFileEntity
	err := result.Decode(&mongoFile)
	if err != nil {
		return domain.File{}, err
	}

	metadata, domainError := mongoFile.ToDomain()
	if domainError != nil {
		return domain.File{}, domainError
	}
	return metadata, nil
}

func (repo *MongoFileRepository) DeleteFile(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	result, err := repo.client.Database(repo.database).Collection(repo.collection).DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount <= 0 {
		return errors.New("could not delete. please try again")
	}
	return nil
}
