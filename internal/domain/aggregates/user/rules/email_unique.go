package rules

import (
	"fmt"
	"physk/internal/domain/aggregates/user/context"
	userErrors "physk/internal/domain/aggregates/user/errors"
	vo "physk/internal/domain/aggregates/user/value_objects"
)

func CheckEmailUniqueness(
	ctx context.UserCtx,
	email vo.UserEmail,
) error {
	exists, err := ctx.Repo.EmailExists(ctx, email)
	if err != nil {
		return fmt.Errorf("check email exists: %w", err)
	}
	if exists {
		return userErrors.EmailNotUniqueError
	}
	return nil
}
