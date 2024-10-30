package util

import (
	"bytes"
	"encoding/base64"
	image2 "image"
	"image/jpeg"
)

func ToJPEG(code string) ([]byte, error) {
	decodedImage, err := base64.StdEncoding.DecodeString(code)
	if err != nil {
		return nil, err
	}

	img, _, err := image2.Decode(bytes.NewReader(decodedImage))
	if err != nil {
		return nil, err
	}

	var imgBuffer bytes.Buffer
	if err = jpeg.Encode(&imgBuffer, img, &jpeg.Options{Quality: 70}); err != nil {
		return nil, err
	}

	jpegData := imgBuffer.Bytes()
	jpegImgCode := make([]byte, base64.StdEncoding.EncodedLen(len(jpegData)))
	base64.StdEncoding.Encode(jpegImgCode, jpegData)

	return jpegImgCode, nil
}
