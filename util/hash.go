package util

import (
	"crypto/sha1"
	"fmt"
)

const salt string = "hjqrhjqw124617ajfhajs"

func GenerateHash(line string) string {
	hash := sha1.New()
	hash.Write([]byte(line))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
