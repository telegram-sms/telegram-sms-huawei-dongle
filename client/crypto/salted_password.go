package crypto

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"strings"
)

func SHA256(text string) string {
	hash := sha256.New()
	hash.Write([]byte(text))
	result := hash.Sum(nil)
	return strings.ToLower(hex.EncodeToString(result))
}

func B64(text string) string {
	blob := []byte(text)
	return base64.StdEncoding.EncodeToString(blob)
}

func b64_sha256(text string) string {
	return B64(SHA256(text))
}

func EncodeSaltedPassword(username, password, token string) string {
	return b64_sha256(username + b64_sha256(password) + token)
}
