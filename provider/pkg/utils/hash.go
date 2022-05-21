package utils

import (
	"crypto/sha1"
	"encoding/base64"
)

func SHA1Hash(s string) string {
	hasher := sha1.New()
	hasher.Write([]byte(s))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
