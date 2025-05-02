package controller

import (
	"context"
	"errors"
	"fmt"
	"physk/internal/domain/aggregates/user"
	"physk/internal/infrastructure/storage"
	tokenServ "physk/internal/services/token_service"
)

func (c *Controller) GetMe(
	ctx context.Context,
	token string,
) (user.User, error) {
	claims, err := c.tokenService.ValidateToken(token)
	if err != nil {
		if !errors.Is(err, tokenServ.InvalidTokenError) {
			return user.User{}, fmt.Errorf("validate token: %w", err)
		}
	}
	usr, err := c.read.GetUserByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, storage.ErrorNotFound) {
			return user.User{}, UserNotFound
		}
		return user.User{}, fmt.Errorf("get user by login: %w", err)
	}
	return usr, nil
}
