package util

import (
	"crypto/sha256"
	"encoding/hex"
)

// GetMD5HashString - generate hash of string
func GetMD5HashString(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	md := hash.Sum(nil)
	return hex.EncodeToString(md)
}
