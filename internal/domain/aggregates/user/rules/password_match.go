package rules

import (
	"fmt"
	userErrors "physk/internal/domain/aggregates/user/errors"
	"physk/pkg/hasher"
)

func ComparePasswordAndHash(new, old string) (newHash string, err error) {
	hash, err := hasher.HashString(new)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}

	if hash != old {
		return "", userErrors.InvalidPasswordError
	}

	return hash, nil
}
