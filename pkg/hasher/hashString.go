package hasher

import (
	"crypto/md5"
	"encoding/base64"
)

func HashString(s string) (string, error) {
	hash := md5.Sum([]byte(s))
	return base64.StdEncoding.EncodeToString(hash[:]), nil
}
