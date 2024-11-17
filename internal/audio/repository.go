package audio

import (
	"ChatApp/internal/audio/models"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repository interface {
	CreateAudio(ctx context.Context, audio models.Audio) error
	GetAudio(ctx context.Context, audioId primitive.ObjectID) ([]byte, error)
}
