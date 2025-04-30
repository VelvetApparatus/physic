package rules

import (
	"fmt"
	"physk/internal/domain/aggregates/user/context"
	userErrors "physk/internal/domain/aggregates/user/errors"
)

func CheckEmailUniqueness(
	ctx context.UserCtx,
	email string,
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
