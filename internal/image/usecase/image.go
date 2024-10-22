package usecase

import (
	"ChatApp/internal/image"
	"ChatApp/internal/image/models"
	"ChatApp/util"
	"context"
	"encoding/base64"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ImageUsecase struct {
	repo image.Repository
}

func NewImageUsecase(repo image.Repository) *ImageUsecase {
	return &ImageUsecase{repo: repo}
}

func (u *ImageUsecase) CreateImage(ctx context.Context, imageCodes []string) ([]primitive.ObjectID, error) {
	var images []models.Image
	for _, code := range imageCodes {
		webpCode, err := util.ToJPEG(code)
		if err != nil {
			return nil, err
		}

		images = append(images, models.Image{
			ID:        primitive.NewObjectID(),
			ImageCode: webpCode,
			CreatedAt: time.Now().Unix(),
		})
	}

	imageIdsInterface, err := u.repo.CreateImage(ctx, images)
	if err != nil {
		return nil, err
	}

	return util.ToType[primitive.ObjectID](imageIdsInterface)
}

func (u *ImageUsecase) GetImage(ctx context.Context, imageId string) ([]byte, error) {
	imageObjectId, err := primitive.ObjectIDFromHex(imageId)
	if err != nil {
		return nil, err
	}

	imageCode, err := u.repo.GetImage(ctx, imageObjectId)
	if err != nil {
		return nil, err
	}

	imageData := make([]byte, base64.StdEncoding.DecodedLen(len(imageCode)))
	n, err := base64.StdEncoding.Decode(imageData, imageCode)
	if err != nil {
		return nil, err
	}

	return imageData[:n], nil
}
