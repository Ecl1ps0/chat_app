package repository

import (
	"ChatApp/internal/audio/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AudioRepository struct {
	db *mongo.Collection
}

func NewAudioRepository(db *mongo.Database, collection string) *AudioRepository {
	return &AudioRepository{db: db.Collection(collection)}
}

func (r *AudioRepository) CreateAudio(ctx context.Context, audio models.Audio) error {
	_, err := r.db.InsertOne(ctx, &audio)

	return err
}

func (r *AudioRepository) GetAudio(ctx context.Context, imageId primitive.ObjectID) ([]byte, error) {
	filter := bson.M{"_id": imageId}

	var audio models.Audio
	if err := r.db.FindOne(ctx, filter).Decode(&audio); err != nil {
		return nil, err
	}

	return audio.Content, nil
}
