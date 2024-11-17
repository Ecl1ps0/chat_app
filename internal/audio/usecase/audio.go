package usecase

import (
	"ChatApp/internal/audio"
	"ChatApp/internal/audio/models"
	"context"
	"encoding/base64"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"mime/multipart"
	"time"
)

type AudioUsecase struct {
	repo audio.Repository
}

func NewAudioUsecase(repo audio.Repository) *AudioUsecase {
	return &AudioUsecase{repo: repo}
}

func (u *AudioUsecase) CreateAudio(ctx context.Context, file multipart.File) (primitive.ObjectID, error) {
	data, err := io.ReadAll(file)
	if err != nil {
		return primitive.NilObjectID, err
	}

	audioBase64 := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(audioBase64, data)

	audioContent := models.Audio{
		ID:        primitive.NewObjectID(),
		Content:   audioBase64,
		CreatedAt: time.Now().Unix(),
	}

	if err = u.repo.CreateAudio(ctx, audioContent); err != nil {
		return primitive.NilObjectID, err
	}

	return audioContent.ID, nil
}

func (u *AudioUsecase) GetAudio(ctx context.Context, audioId string) ([]byte, error) {
	hex, err := primitive.ObjectIDFromHex(audioId)
	if err != nil {
		return nil, err
	}

	audioCode, err := u.repo.GetAudio(ctx, hex)
	if err != nil {
		return nil, err
	}

	audioData := make([]byte, base64.StdEncoding.DecodedLen(len(audioCode)))
	n, err := base64.StdEncoding.Decode(audioData, audioCode)
	if err != nil {
		return nil, err
	}

	return audioData[:n], nil
}
