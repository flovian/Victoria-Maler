package services

import "crypto/sha256"

func Hash(data string) string {
    h := sha256.Sum256([]byte(data))
    return string(h[:])
}
