package utils

import (
	"crypto/rand"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

func GenerateToken(userID uint) (string, error) {
	b := make([]byte, 16)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	token := hex.EncodeToString(b)

	return token, nil
}

func GenerateFromPasswordString(p string) string {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)

	return string(hashed)
}

func GenerateFromPasswordByte(p string) []byte {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)

	return hashed
}
