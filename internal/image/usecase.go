package image

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Usecase interface {
	CreateImage(ctx context.Context, images []string) ([]primitive.ObjectID, error)
	GetImage(ctx context.Context, imageId string) ([]byte, error)
}
