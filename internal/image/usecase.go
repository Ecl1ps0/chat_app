package image

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mime/multipart"
)

type Usecase interface {
	CreateImages(ctx context.Context, files []multipart.File) ([]primitive.ObjectID, error)
	CreateImage(ctx context.Context, file multipart.File) (primitive.ObjectID, error)
	GetImage(ctx context.Context, imageId string) ([]byte, error)
}
