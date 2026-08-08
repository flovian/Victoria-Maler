package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func Sum(data string) []byte {
	h := sha256.Sum256([]byte(data))
	return h[:]
}

func HashHex(data string) string {
	return hex.EncodeToString(Sum(data))
}

func HashHexBytes(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}
