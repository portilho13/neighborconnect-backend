package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"log"
)

func GenerateRandomHash() string {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Fatal("Enable to generate random Hash")
	}

	hash := sha256.Sum256(randomBytes)
	return hex.EncodeToString(hash[:])

}