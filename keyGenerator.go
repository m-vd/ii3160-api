package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

//GenerateAPIKey is a function that provides a key
func GenerateAPIKey() string {
	key := [256]byte{}
	_, err := rand.Read(key[:])
	if err != nil {
		panic(err)
	}
	apiKey := sha256.Sum256(key[:])

	return hex.EncodeToString(apiKey[:])
}
