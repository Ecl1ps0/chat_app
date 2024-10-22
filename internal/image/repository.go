package image

import (
	"ChatApp/internal/image/models"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repository interface {
	CreateImage(ctx context.Context, images []models.Image) ([]interface{}, error)
	GetImage(ctx context.Context, imageId primitive.ObjectID) ([]byte, error)
}
