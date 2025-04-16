package usecase

import (
	"ChatApp/internal/image"
	"ChatApp/internal/image/models"
	"ChatApp/util"
	"context"
	"encoding/base64"
	"errors"
	"github.com/sunshineplan/imgconv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	image2 "image"
	"mime/multipart"
	"time"
)

type ImageUsecase struct {
	repo image.Repository
}

func NewImageUsecase(repo image.Repository) *ImageUsecase {
	return &ImageUsecase{repo: repo}
}

func (u *ImageUsecase) CreateImages(ctx context.Context, files []multipart.File) ([]primitive.ObjectID, error) {
	var imageCodes []models.Image
	for _, file := range files {
		img, err := imgconv.Decode(file)
		if err != nil {
			return nil, err
		}

		jpegCode, err := util.ToJPEGBase64(img)
		if err != nil {
			return nil, err
		}

		imageCodes = append(imageCodes, models.Image{
			ID:        primitive.NewObjectID(),
			ImageCode: jpegCode,
			CreatedAt: time.Now().Unix(),
		})
	}

	imageIdsInterface, err := u.repo.CreateImages(ctx, imageCodes)
	if err != nil {
		return nil, err
	}

	return util.ToType[primitive.ObjectID](imageIdsInterface)
}

func (u *ImageUsecase) CreateImage(ctx context.Context, file multipart.File) (primitive.ObjectID, error) {
	decodedImage, _, err := image2.Decode(file)
	if err != nil {
		return primitive.NilObjectID, err
	}

	jpegCode, err := util.ToJPEGBase64(decodedImage)
	if err != nil {
		return primitive.NilObjectID, err
	}

	img := models.Image{
		ID:        primitive.NewObjectID(),
		ImageCode: jpegCode,
		CreatedAt: time.Now().Unix(),
	}

	imgIdInterface, err := u.repo.CreateImage(ctx, img)
	if err != nil {
		return primitive.NilObjectID, err
	}

	imgId, ok := imgIdInterface.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, errors.New("Fail to parse image id from interface")
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
