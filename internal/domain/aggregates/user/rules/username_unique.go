package rules

import (
	"fmt"
	"physk/internal/domain/aggregates/user/context"
	userErrors "physk/internal/domain/aggregates/user/errors"
	vo "physk/internal/domain/aggregates/user/value_objects"
)

func CheckUsernameUniqueness(
	userCtx context.UserCtx,
	username vo.UserName,
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
