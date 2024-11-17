package audio

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mime/multipart"
)

type Usecase interface {
	CreateAudio(ctx context.Context, file multipart.File) (primitive.ObjectID, error)
	GetAudio(ctx context.Context, imageId string) ([]byte, error)
}
