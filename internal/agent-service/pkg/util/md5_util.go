package util

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(params string) string {
	hash := md5.New()

	// Write the data to the hash
	hash.Write([]byte(params))

	// Get the resulting hash as a byte slice
	hashBytes := hash.Sum(nil)

	// Convert to hexadecimal string
	return hex.EncodeToString(hashBytes)
}
