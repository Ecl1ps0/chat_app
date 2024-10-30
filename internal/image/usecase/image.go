package usecase

import (
	"ChatApp/internal/image"
	"ChatApp/internal/image/models"
	"ChatApp/util"
	"context"
	"encoding/base64"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ImageUsecase struct {
	repo image.Repository
}

func NewImageUsecase(repo image.Repository) *ImageUsecase {
	return &ImageUsecase{repo: repo}
}

func (u *ImageUsecase) CreateImages(ctx context.Context, imageCodes []string) ([]primitive.ObjectID, error) {
	var images []models.Image
	for _, code := range imageCodes {
		jpegCode, err := util.ToJPEG(code)
		if err != nil {
			return nil, err
		}

		images = append(images, models.Image{
			ID:        primitive.NewObjectID(),
			ImageCode: jpegCode,
			CreatedAt: time.Now().Unix(),
		})
	}

	imageIdsInterface, err := u.repo.CreateImages(ctx, images)
	if err != nil {
		return nil, err
	}

	return util.ToType[primitive.ObjectID](imageIdsInterface)
}

func (u *ImageUsecase) CreateImage(ctx context.Context, imageCode string) (primitive.ObjectID, error) {
	jpegCode, err := util.ToJPEG(imageCode)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	img := models.Image{
		ID:        primitive.NewObjectID(),
		ImageCode: jpegCode,
		CreatedAt: time.Now().Unix(),
	}

	imgIdInterface, err := u.repo.CreateImage(ctx, img)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	imgId, ok := imgIdInterface.(primitive.ObjectID)
	if !ok {
		return primitive.ObjectID{}, errors.New("Fail to parse image id from interface")
	}

	return imgId, nil
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
