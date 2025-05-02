package controller

import (
	"context"
	"errors"
	"fmt"
	"physk/internal/domain/aggregates/user"
	uCtx "physk/internal/domain/aggregates/user/context"
	deliveryDTO "physk/internal/infrastructure/delivery/dto"
	"physk/internal/infrastructure/storage"
)

var (
	UserNotFound = errors.New("user not found")
)

func (c *Controller) RegisterUser(
	ctx context.Context,
	request deliveryDTO.UserCreateRequest,
) (user.User, error) {
	dto := request.ToAggregateDTO()
	userCtx := uCtx.NewUserContext(ctx)

	userAgg, err := user.CreateUser(userCtx, dto)
	if err != nil {
		return user.User{}, fmt.Errorf("create user: %w", err)
	}

	err = c.write.Create(ctx, userAgg)
	if err != nil {
		return user.User{}, fmt.Errorf("create in write model: %w", err)
	}

	return userAgg, nil
}

func (c *Controller) Login(
	ctx context.Context,
	request deliveryDTO.UserLoginRequest,
) (string, error) {
	dto := request.ToAggregateDTO()
	userCtx := uCtx.NewUserContext(ctx)

	userAgg, err := c.read.GetUserByLogin(userCtx, dto.Login)
	if err != nil {
		if errors.Is(err, storage.ErrorNotFound) {
			return "", UserNotFound
		}
		return "", fmt.Errorf("get user by login: %w", err)
	}

	err = userAgg.MatchPassword(request.Password)
	if err != nil {
		return "", fmt.Errorf("match password: %w", err)
	}

	token, err := c.tokenService.NewToken(userAgg)
	if err != nil {
		return "", fmt.Errorf("gen new token: %w", err)
	}

	return token, nil
}
