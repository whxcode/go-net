package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateToken(userID uint) (string, error) {
	b := make([]byte, 16)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	token := hex.EncodeToString(b)

	return token, nil
}
