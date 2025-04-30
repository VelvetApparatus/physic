package rules

import (
	"fmt"
	"physk/internal/domain/aggregates/user/context"
	userErrors "physk/internal/domain/aggregates/user/errors"
)

func CheckUsernameUniqueness(
	userCtx context.UserCtx,
	username string,
) error {
	exists, err := userCtx.Repo.UsernameExists(userCtx, username)
	if err != nil {
		return fmt.Errorf("check username exists: %w", err)
	}

	if exists {
		return userErrors.UsernameNotUniqueError
	}
	return nil
}
