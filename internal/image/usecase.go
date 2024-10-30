package image

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Usecase interface {
	CreateImages(ctx context.Context, imageCodes []string) ([]primitive.ObjectID, error)
	CreateImage(ctx context.Context, imageCode string) (primitive.ObjectID, error)
	GetImage(ctx context.Context, imageId string) ([]byte, error)
}
