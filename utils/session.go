package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateSessionID() (string, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to generate session id %w", err)
	}

	return base64.URLEncoding.EncodeToString(b), nil
}
