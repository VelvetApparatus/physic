package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"physk/internal/domain/aggregates/user"
	"physk/internal/infrastructure/storage"
)

func (c *Controller) GetMe(
	ctx context.Context,
	userID uuid.UUID,
) (user.User, error) {
	usr, err := c.read.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrorNotFound) {
			return user.User{}, UserNotFound
		}
		return user.User{}, fmt.Errorf("get user by login: %w", err)
	}
	return usr, nil
}
