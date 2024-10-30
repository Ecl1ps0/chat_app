package repository

import (
	"ChatApp/internal/image/models"
	"ChatApp/util"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ImageRepository struct {
	db *mongo.Collection
}

func NewImageRepository(db *mongo.Database, collection string) *ImageRepository {
	return &ImageRepository{db: db.Collection(collection)}
}

func (r *ImageRepository) CreateImages(ctx context.Context, images []models.Image) ([]interface{}, error) {
	imageInterfaces := util.ToInterfaceSlice(images)

	result, err := r.db.InsertMany(ctx, imageInterfaces)
	if err != nil {
		return nil, err
	}

	return result.InsertedIDs, nil
}

func (r *ImageRepository) CreateImage(ctx context.Context, image models.Image) (interface{}, error) {
	result, err := r.db.InsertOne(ctx, image)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	return result.InsertedID, nil
}

func (r *ImageRepository) GetImage(ctx context.Context, imageId primitive.ObjectID) ([]byte, error) {
	filter := bson.M{"_id": imageId}

	var image models.Image
	if err := r.db.FindOne(ctx, filter).Decode(&image); err != nil {
		return nil, err
	}

	return image.ImageCode, nil
}
