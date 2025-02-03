package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func GenerateSalt() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

func EncryptPassword(password string, salt string) string {
	h := sha256.New()
	h.Write([]byte(password + salt))
	return hex.EncodeToString(h.Sum(nil))
}

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateToken(length uint) string {
	if length == 0 {
		return ""
	}

	tokenBytes := make([]byte, length)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		panic("Failed to generate random bytes: " + err.Error())
	}

	var builder strings.Builder
	for i := range tokenBytes {
		index := int(tokenBytes[i]) % len(alphabet)
		builder.WriteByte(alphabet[index])
	}

	return builder.String()
}
