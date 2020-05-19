package main

import (
	"crypto/sha256"
	"encoding/hex"
)

func hash(value string) string {
	if value == "" {
		return ""
	}
	checksum := sha256.Sum256([]byte(value))
	asBytes := checksum[:]
	asStr := hex.EncodeToString(asBytes)
	return asStr
}
