package hasher

import (
	"crypto/md5"
	"fmt"
)

func HashString(s string) (string, error) {
	md5Hash := md5.New()
	_, err := md5Hash.Write([]byte(s))
	if err != nil {
		return "", fmt.Errorf("cannot write string to hasher writer: %w", err)
	}
	return string(md5Hash.Sum(nil)), nil
}
