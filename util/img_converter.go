package util

import (
	"bytes"
	"encoding/base64"
	image2 "image"
	"image/jpeg"
)

func ToJPEGBase64(image image2.Image) ([]byte, error) {
	var jpegData bytes.Buffer
	if err := jpeg.Encode(&jpegData, image, &jpeg.Options{Quality: 70}); err != nil {
		return nil, err
	}

	jpegImgCode := make([]byte, base64.StdEncoding.EncodedLen(len(jpegData.Bytes())))
	base64.StdEncoding.Encode(jpegImgCode, jpegData.Bytes())

	return jpegImgCode, nil
}
