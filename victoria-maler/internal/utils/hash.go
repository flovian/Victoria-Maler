package utils

import "crypto/sha256"

func Sum(data string) []byte {
    h := sha256.Sum256([]byte(data))
    return h[:]
}
