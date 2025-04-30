package rules

import (
	"fmt"
	"physk/internal/domain/aggregates/user/context"
	userErrors "physk/internal/domain/aggregates/user/errors"
)

func CheckLoginUniqueness(
	ctx context.UserCtx,
	login string,
) error {
	exists, err := ctx.Repo.LoginExists(ctx, login)
	if err != nil {
		return fmt.Errorf("check login exists: %w", err)
	}

	if exists {
		return userErrors.LoginNotUniqueError
	}
	return nil
}
