package util

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/pkg/errors"
)

func GenerateRandomDataString(length int) (string, error) {
	randData := make([]byte, length)
	if _, err := rand.Read(randData); err != nil {
		return "", errors.WithStack(err)
	}
	return base64.URLEncoding.EncodeToString(randData), nil
}
