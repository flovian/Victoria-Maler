package services

import (
	"crypto/sha256"
	"encoding/hex"
)

// Hash returns the hex-encoded SHA-256 digest of the input.
// This matches the format used across the platform for Bitcoin anchoring.
func Hash(data string) string {
	h := sha256.Sum256([]byte(data))
	return hex.EncodeToString(h[:])
}
